package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/url"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

var ErrInvalidSignature = errors.New("session: signature not verified")

func generateSecret(base, salt string, keySize int) []byte {
	return pbkdf2.Key([]byte(base), []byte(salt), keyIterNum, keySize, sha1.New)
}

// The origin of this snippet can be found at https://gist.github.com/doitian/2a89dc9e4372e55335c9111f576b47bf
func verifySign(encryptedData, sign, base, signSalt string) (bool, error) {
	signKey := generateSecret(base, signSalt, keySize)
	signHmac := hmac.New(sha1.New, signKey)
	signHmac.Write([]byte(encryptedData))
	verifySign := signHmac.Sum(nil)
	signDecoded, err := hex.DecodeString(sign)
	if err != nil {
		return false, err
	}
	if !hmac.Equal(verifySign, signDecoded) {
		return false, ErrInvalidSignature
	}
	return true, nil
}

func decodeCookieData(cookie []byte) (data, iv []byte, err error) {
	vectors := strings.SplitN(string(cookie), "--", 2)

	data, err = base64.StdEncoding.DecodeString(vectors[0])
	if err != nil {
		return
	}

	iv, err = base64.StdEncoding.DecodeString(vectors[1])
	if err != nil {
		return
	}

	return
}

func decryptCookie(cookie []byte, secret []byte) (dd []byte, err error) {
	data, iv, err := decodeCookieData(cookie)

	c, err := aes.NewCipher(secret[:32])
	if err != nil {
		return
	}

	cfb := cipher.NewCBCDecrypter(c, iv)
	dd = make([]byte, len(data))
	cfb.CryptBlocks(dd, data)

	return
}

func decryptCookieWithAuthentication(data, iv, authTag, secret []byte) (dd []byte, err error) {
	c, err := aes.NewCipher(secret[:32])
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return
	}
	if gcm.Overhead() != len(authTag) {
		return nil, ErrInvalidSignature
	}
	dd, err = gcm.Open(nil, iv, append(data, authTag...), nil)
	if err != nil {
		return nil, ErrInvalidSignature
	}

	return
}

func DecryptSignedCookie(signedCookie, secretKeyBase, salt, signSalt string) (session []byte, err error) {
	cookie, err := url.QueryUnescape(signedCookie)
	if err != nil {
		return
	}

	vectors := strings.SplitN(cookie, "--", 2)
	if len(vectors) != 2 || vectors[0] == "" || vectors[1] == "" {
		return nil, errors.New("session: invalid cookie")
	}
	verified, err := verifySign(vectors[0], vectors[1], secretKeyBase, signSalt)
	if err != nil {
		return
	}
	if !verified {
		return nil, ErrInvalidSignature
	}

	data, err := base64.StdEncoding.DecodeString(vectors[0])
	if err != nil {
		return
	}

	session, err = decryptCookie(data, generateSecret(secretKeyBase, salt, keySize))
	if err != nil {
		return
	}

	return
}

func DecryptAuthorizedEncryptedCookie(signedCookie, secretKeyBase, salt string) (session []byte, err error) {
	cookie, err := url.QueryUnescape(signedCookie)
	if err != nil {
		return
	}

	vectors := strings.SplitN(cookie, "--", 3)
	if len(vectors) != 3 || vectors[0] == "" || vectors[1] == "" || vectors[2] == "" {
		return nil, errors.New("session: invalid cookie")
	}

	data, err := base64.StdEncoding.DecodeString(vectors[0])
	if err != nil {
		return
	}

	iv, err := base64.StdEncoding.DecodeString(vectors[1])
	if err != nil {
		return
	}

	authTag, err := base64.StdEncoding.DecodeString(vectors[2])
	if err != nil {
		return
	}

	session, err = decryptCookieWithAuthentication(data, iv, authTag, generateSecret(secretKeyBase, salt, keySizeGCM))
	if err != nil {
		return
	}

	return
}

// Rails 4.0 defaults
const (
	keyIterNum = 1000
	keySize    = 64
	keySizeGCM = 32
)

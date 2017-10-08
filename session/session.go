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

func generateSecret(base, salt string) []byte {
	return pbkdf2.Key([]byte(base), []byte(salt), keyIterNum, keySize, sha1.New)
}

// I stolen the snippet from: https://gist.github.com/doitian/2a89dc9e4372e55335c9111f576b47bf
func verifySign(encryptedData, sign, base, signSalt string) bool {
	signKey := generateSecret(base, signSalt)
	signHmac := hmac.New(sha1.New, signKey)
	signHmac.Write([]byte(encryptedData))
	verifySign := signHmac.Sum(nil)
	signDecoded, err := hex.DecodeString(sign)
	if err != nil {
		return false
	}
	if hmac.Equal(verifySign, signDecoded) {
		return true
	}
	return false
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

func DecryptSignedCookie(signedCookie, secretKeyBase, salt, signSalt string) (session []byte, err error) {
	cookie, err := url.QueryUnescape(signedCookie)
	if err != nil {
		return
	}

	vectors := strings.SplitN(cookie, "--", 2)
	if vectors[0] == "" || vectors[1] == "" {
		return nil, errors.New("Invalid cookie!")
	}
	if verified := verifySign(vectors[0], vectors[1], secretKeyBase, signSalt); !verified {
		return nil, errors.New("Not verified signature!")
	}

	data, err := base64.StdEncoding.DecodeString(vectors[0])
	if err != nil {
		return
	}

	session, err = decryptCookie(data, generateSecret(secretKeyBase, salt))
	if err != nil {
		return
	}

	return
}

// Rails 4.0 defaults
const (
	keyIterNum = 1000
	keySize    = 64
)

package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/url"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

var ErrInvalidSignature = errors.New("session: signature not verified")

func generateSecret(base, salt string) []byte {
	return pbkdf2.Key([]byte(base), []byte(salt), keyIterNum, keySize, sha1.New)
}

func signData(encryptedData, base, signSalt string) []byte {
	signKey := generateSecret(base, signSalt)
	signHmac := hmac.New(sha1.New, signKey)
	signHmac.Write([]byte(encryptedData))
	return signHmac.Sum(nil)
}

// The origin of this snippet can be found at https://gist.github.com/doitian/2a89dc9e4372e55335c9111f576b47bf
func verifySign(encryptedData, sign, base, signSalt string) (bool, error) {
	verifySign := signData(encryptedData, base, signSalt)
	signDecoded, err := hex.DecodeString(sign)
	if err != nil {
		return false, err
	}
	if !hmac.Equal(verifySign, signDecoded) {
		return false, ErrInvalidSignature
	}
	return true, nil
}

// sign and join data with signature using "--" (needs to be url.QueryEscape'd)
func signJoiner(encryptedData, base, signSalt string) string {
	postfix := hex.EncodeToString(signData(encryptedData, base, signSalt))
	return strings.Join([]string{encryptedData, postfix}, "--")
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

func encodeCookieData(data, iv []byte) (cookie []byte) {
	datas := base64.StdEncoding.EncodeToString(data)
	ivs := base64.StdEncoding.EncodeToString(iv)
	cookie = []byte(strings.Join([]string{datas, ivs}, "--"))
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

// padSession implements PKCS#7 padding for the plaintext
// https://en.wikipedia.org/wiki/Padding_(cryptography)#PKCS7
func padSession(session []byte, blockSize int) []byte {
	sesslen := len(session)
	padsize := blockSize - (sesslen % blockSize)
	if padsize == blockSize {
		return session
	}
	newlen := sesslen + padsize
	padbyte := byte(padsize)
	padded := make([]byte, newlen)
	copy(padded, session)
	for i := sesslen; i < newlen; i++ {
		padded[i] = padbyte
	}

	return padded
}

func encryptCookie(dd, secret []byte) (cookie []byte, err error) {
	c, err := aes.NewCipher(secret[:32])
	if err != nil {
		return
	}
	padded := padSession(dd, c.BlockSize())
	iv := make([]byte, c.BlockSize())
	// rails uses a random iv, so this should be fine: https://github.com/rails/rails/blob/master/activesupport/lib/active_support/message_encryptor.rb#L172
	_, err = rand.Read(iv)
	if err != nil {
		return
	}

	cfb := cipher.NewCBCEncrypter(c, iv)
	data := make([]byte, len(padded))
	cfb.CryptBlocks(data, padded)

	cookie = encodeCookieData(data, iv)

	return
}

// EncryptSignedCookie encrypts and signs session to produce a cookie that rails can read
func EncryptSignedCookie(session []byte, secretKeyBase, salt, signSalt string) (signedCookie string, err error) {
	data, err := encryptCookie(session, generateSecret(secretKeyBase, salt))
	if err != nil {
		return
	}

	datastr := base64.StdEncoding.EncodeToString(data)
	cookie := signJoiner(datastr, secretKeyBase, signSalt)
	signedCookie = url.QueryEscape(cookie)
	return
}

func DecryptSignedCookie(signedCookie, secretKeyBase, salt, signSalt string) (session []byte, err error) {
	cookie, err := url.QueryUnescape(signedCookie)
	if err != nil {
		return
	}

	vectors := strings.SplitN(cookie, "--", 2)
	if vectors[0] == "" || vectors[1] == "" {
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

package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"net/url"
	"strings"

	"code.google.com/p/go.crypto/pbkdf2"
)

func generateSecret(base, salt string) []byte {
	return pbkdf2.Key([]byte(base), []byte(salt), key_iter_num, key_size, sha1.New)
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

func DecryptSignedCookie(signed_cookie, secret_key_base, salt string) (session []byte, err error) {
	cookie, err := url.QueryUnescape(signed_cookie)
	if err != nil {
		return
	}

	vectors := strings.SplitN(cookie, "--", 2)
	data, err := base64.StdEncoding.DecodeString(vectors[0])
	if err != nil {
		return
	}

	session, err = decryptCookie(data, generateSecret(secret_key_base, salt))
	if err != nil {
		return
	}

	return
}

// Rails 4.0 defaults
const (
	key_iter_num = 1000
	key_size     = 64
)

package session

import (
	"bytes"
	"net/url"
	"strings"
	"testing"
)

const (
	secretKeyBase = "fe98c394d54eeae9edff39c1934b156607e4376188463d397d460eef9585cf15c0dd23f353877552d1c9b0565a03b7fdeadfb33907c6d582eb02319a7409610b"
	// The cookie's original content is:
	// map[flash:map[discard:[] flashes:map[notice:Welcome! You have signed up successfully.]] session_id:b85897340bfedc7e03b7e9479c271439 _csrf_token:dTDcQiGuEE8n6KUQmXNhIoXsLQJlqrBPUAsspGMpkdg= warden.user.user.key:[[1] $2a$11$6omJ7/e3Ni7Pl7jZbCdDBu]]
	signedCookie = "RkpiOStFLzExVm42aXZiMFZWaDB3c09rbEE4aTUvcEg5Q1VnaTNDOTBwMTdSUGFsdjZqbWZpQmV3eXhQbEJieE1EYXZCQXNGNFhKREI5aUx0aXVFZE1vaXQzSTdtYzc5S1NmeXBEZG93Mm1PQmQ2RVMvdjRqbTdsTW1qTjcxRTZFSVpCZFBUcTByN0ZYQmhWWVZPVE45RUsyS2NRcEV5QkdsajRUL3FGYjNmdUZrYmZ5TVZxSlpucllOaXlTN0pZZG85eHlMNEN0MVdYayttdE8wNTBTSElDYTRqditGMmpoL09hcDhkTFZ0dngyM244aG53aWNLNWRvVTN3K2dpUWd0eGttRXZUdGx2TGJHS0xlN0hKWFI2aVhuQlE4Y3NvYWx1QTZvcDRkbDJZdjl4NGJ1b1B1WW9QdXdEOVpzcCtBR1BCVDkxZkNSVENJZkVqMkgzR3pxQ1lVVEJmQlBYK0ZIQWJ5WHRpOC84PS0taDluekdrZE1LbzVrZDVlMHFSSzNjdz09--5f676b46cb0671630fd33bfec08b6fbf3f858c6a"
	salt         = "encrypted cookie"
	signSalt     = "signed encrypted cookie"

	// The cookie's original content is:
	plainjson = `{"session_id":"b85897340bfedc7e03b7e9479c271439","_csrf_token":"dTDcQiGuEE8n6KUQmXNhIoXsLQJlqrBPUAsspGMpkdg=","warden.user.user.key":[[1],"$2a$11$6omJ7/e3Ni7Pl7jZbCdDBu"],"flash":{"discard":[],"flashes":{"notice":"Welcome! You have signed up successfully."}}}`
)

type sessionObj struct {
	SessionID string `json:"session_id"`
	CSRF      string `json:"_csrf_token"`
}

func TestVerifySign(t *testing.T) {
	cookie, _ := url.QueryUnescape(signedCookie)
	vectors := strings.SplitN(cookie, "--", 2)
	// a valid signature case
	verified, err := verifySign(vectors[0], vectors[1], secretKeyBase, signSalt)
	if err != nil {
		t.Errorf("verifySign test failure: %v", err)
	}
	if !verified {
		t.Errorf("verifySign test failure: %v", ErrInvalidSignature)
	}
	// an invalid signature case
	faultSignSalt := "wrong signature salt"
	verified, err = verifySign(vectors[0], vectors[1], secretKeyBase, faultSignSalt)
	if err == nil || verified {
		t.Error("verifySign test with an invalid signature salt passed")
	}
}

func TestSignJoiner(t *testing.T) {
	encryptedData := "RkpiOStFLzExVm42aXZiMFZWaDB3c09rbEE4aTUvcEg5Q1VnaTNDOTBwMTdSUGFsdjZqbWZpQmV3eXhQbEJieE1EYXZCQXNGNFhKREI5aUx0aXVFZE1vaXQzSTdtYzc5S1NmeXBEZG93Mm1PQmQ2RVMvdjRqbTdsTW1qTjcxRTZFSVpCZFBUcTByN0ZYQmhWWVZPVE45RUsyS2NRcEV5QkdsajRUL3FGYjNmdUZrYmZ5TVZxSlpucllOaXlTN0pZZG85eHlMNEN0MVdYayttdE8wNTBTSElDYTRqditGMmpoL09hcDhkTFZ0dngyM244aG53aWNLNWRvVTN3K2dpUWd0eGttRXZUdGx2TGJHS0xlN0hKWFI2aVhuQlE4Y3NvYWx1QTZvcDRkbDJZdjl4NGJ1b1B1WW9QdXdEOVpzcCtBR1BCVDkxZkNSVENJZkVqMkgzR3pxQ1lVVEJmQlBYK0ZIQWJ5WHRpOC84PS0taDluekdrZE1LbzVrZDVlMHFSSzNjdz09"
	if want, got := signedCookie, url.QueryEscape(signJoiner(encryptedData, secretKeyBase, signSalt)); want != got {
		t.Errorf("expected: %q\ngot: %q", want, got)
	}
}

func TestEncryptCookie(t *testing.T) {
	// get plaintext
	plaintext, _ := DecryptSignedCookie(signedCookie, secretKeyBase, salt, signSalt)
	ciphertext, _ := encryptCookie(plaintext, generateSecret(secretKeyBase, salt))
	plaintext2, _ := decryptCookie(ciphertext, generateSecret(secretKeyBase, salt))
	if !bytes.Equal(plaintext, plaintext2) {
		t.Errorf("ciphertext output by encryptCookie cannot be decrypted by decryptCookie")
	}
}

func TestEncryptSignedCookie(t *testing.T) {
	want := make([]byte, len(plainjson)+13)
	copy(want, plainjson)
	for i := len(plainjson); i < len(plainjson)+13; i++ {
		want[i] = 13
	}

	ciphertext, _ := EncryptSignedCookie([]byte(plainjson), secretKeyBase, salt, signSalt)
	got, err := DecryptSignedCookie(ciphertext, secretKeyBase, salt, signSalt)
	if err != nil {
		t.Errorf("got error from decrypting ciphertext: %v", err)
	}
	if !bytes.Equal(want, got) {
		t.Errorf("decrypted ciphertext %q does not match plaintext %q", string(got), string(want))
	}
}

func TestDecryptSignedCookie(t *testing.T) {
	want := make([]byte, len(plainjson)+13)
	copy(want, plainjson)
	for i := len(plainjson); i < len(plainjson)+13; i++ {
		want[i] = 13
	}

	got, err := DecryptSignedCookie(signedCookie, secretKeyBase, salt, signSalt)
	if err != nil {
		t.Errorf("DecryptSignedCookie test failure: %v", err)
	}
	if !bytes.Equal(want, got) {
		t.Errorf("decrypted ciphertext %q does not match plaintext %q", string(got), string(want))
	}
}

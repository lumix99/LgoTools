package jwt

import (
	"crypto"
	"crypto/hmac"
	_ "crypto/sha256"
	"errors"
)

type Signature interface {
	Verify(string, string, string) error
	Sign(string, string) (string, error)
}

type SignByHMAC struct {
	Alg  string
	Hash crypto.Hash
}

var SIGNSHA256 *SignByHMAC

func init() {
	SIGNSHA256 = &SignByHMAC{Alg: "HS256", Hash: crypto.SHA256}
}

func getSignature(alg string) Signature {

	switch alg {
	default:
		return SIGNSHA256
	}
}

func (s *SignByHMAC) Verify(base64HP, sign, secret string) error {

	if !s.Hash.Available() {
		return errors.New("hash function is unavailable")
	}

	signBytes, err := Base64Decode(sign)
	if err != nil {
		return err
	}

	hasher := hmac.New(s.Hash.New, []byte(secret))
	hasher.Write([]byte(base64HP))

	if !hmac.Equal(signBytes, hasher.Sum(nil)) {
		return errors.New("invalid sign")
	}

	return nil

}

func (s *SignByHMAC) Sign(base64HP, secret string) (string, error) {

	keyBytes := []byte(secret)

	if !s.Hash.Available() {
		return "", errors.New("hash function is unavailable")
	}

	hasher := hmac.New(s.Hash.New, keyBytes)
	hasher.Write([]byte(base64HP))

	signStr := Base64Encode(hasher.Sum(nil))

	return signStr, nil
}

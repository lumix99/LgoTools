package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

type Header struct {
	Htype string `json:"type"`
	Alg   string `json:"alg"`
}

func NewHeader(alg string) Header {
	switch alg {
	case "HS256":
	default:
		alg = "HS256"
	}

	return Header{Htype: "JWT", Alg: alg}
}

func Create(header Header, payload Payload, secret string) (string, error) {

	sl := make([]string, 2)

	if jsonVal, err := json.Marshal(header); err != nil {
		return "", err
	} else {
		sl[0] = Base64Encode(jsonVal)
	}

	if jsonVal, err := json.Marshal(payload); err != nil {
		return "", err
	} else {
		sl[1] = Base64Encode(jsonVal)
	}

	base64HP := strings.Join(sl, ".")

	sgt := getSignature(header.Alg)

	sign, err := sgt.Sign(base64HP, secret)
	if err != nil {
		return "", err
	}

	return strings.Join([]string{base64HP, sign}, "."), nil
}

func Verify(token, secret string, payload Payload) error {

	parts := strings.Split(token, ".")

	if len(parts) != 3 {
		return errors.New("invalid token")
	}

	var err error
	header := new(Header)
	var headerBytes []byte

	if headerBytes, err = Base64Decode(parts[0]); err != nil {
		return err
	}

	if err = json.Unmarshal(headerBytes, header); err != nil {
		return err
	}

	var payloadBytes []byte

	if payloadBytes, err = Base64Decode(parts[1]); err != nil {
		return err
	}

	if err = json.Unmarshal(payloadBytes, payload); err != nil {
		return err
	}

	if err = payload.Valid(); err != nil {
		return err
	}

	sgt := getSignature(header.Alg)

	if err = sgt.Verify(strings.Join(parts[0:2], "."), parts[2], secret); err != nil {
		return err
	}

	return nil
}

func Base64Encode(seg []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(seg), "=")
}

func Base64Decode(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}

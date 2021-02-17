package jwt

import (
	"errors"
	"time"
)

const DF_EFFECT_TIME time.Duration = 30 * time.Minute

type Payload interface {
	Valid() error
}

type MapPayload map[string]interface{}

type DefaultPayload struct {
	Iss string `json:"iss,omitempty"` //(issuer)：签发人
	Exp int64  `json:"exp,omitempty"` //(expiration time)：过期时间
	Sub string `json:"sub,omitempty"` //(subject)：主题
	Aud string `json:"aud,omitempty"` //(audience)：受众
	Nbf int64  `json:"nbf,omitempty"` //(Not Before)：生效时间
	Iat int64  `json:"iat,omitempty"` //(Issued At)：签发时间
	Jti string `json:"jti,omitempty"` //(JWT ID)：编号
}

func (m MapPayload) Valid() error {
	return nil
}

func (dp *DefaultPayload) Valid() error {
	now := time.Now().Unix()

	if dp.Exp != 0 && dp.Exp < now {
		return errors.New("jwt is expired")
	}

	if dp.Nbf != 0 && dp.Nbf > now {
		return errors.New("jwt is not valid yet")
	}

	if dp.Iat != 0 && dp.Iat > now {
		return errors.New("invalid Ita")
	}

	return nil
}

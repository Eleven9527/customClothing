package token

import (
	"customClothing/src/config"
	errors "customClothing/src/error"
	"customClothing/src/utils/aes"
	"customClothing/src/utils/base64"
)

func EncodeToken(phone, displayName string) ([]byte, errors.Error) {
	info := phone + config.Cfg().TokenCfg.TokenDelimiter + displayName //phone-displayName

	token, err := aes.EncryptCBC([]byte(info), []byte(config.Cfg().TokenCfg.EncryptKey))
	if err != nil {
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return base64.Encode(token), nil
}

//func DecodeToken(token string) (string, string) {
//
//}

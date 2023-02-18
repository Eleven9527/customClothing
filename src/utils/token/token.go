package token

import (
	"customClothing/src/config"
	errors "customClothing/src/error"
	"customClothing/src/utils/aes"
	"customClothing/src/utils/base64"
	"strings"
)

func EncodeToken(phone, displayName string) ([]byte, errors.Error) {
	info := phone + config.Cfg().TokenCfg.TokenDelimiter + displayName //phone-displayName

	token, err := aes.EncryptCBC([]byte(info), []byte(config.Cfg().TokenCfg.EncryptKey))
	if err != nil {
		return nil, errors.New(errors.INTERNAL_ERROR, "")
	}

	return base64.Encode(token), nil
}

// DecodeToken 解码token，返回phone
func DecodeToken(token []byte) (string, error) {
	tk, err := base64.Decode(token)
	if err != nil {
		return "", err
	}

	//tk = phone-displayName
	tk, err = aes.DecryptCBC(tk, []byte(config.Cfg().TokenCfg.EncryptKey))
	if err != nil {
		return "", err
	}

	infos := strings.Split(string(tk), config.Cfg().TokenCfg.TokenDelimiter)
	return infos[0], nil
}

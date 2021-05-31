package service

import (
	"encoding/base64"
)

func EncodeForRsaBase64String(baseByte []byte) (string, error) {
	rsaByte, err := RsaEncrypt(baseByte)
	if err != nil {
		return "", err
	}
	base64Rsa := base64.StdEncoding.EncodeToString(rsaByte)
	return base64Rsa, nil
}

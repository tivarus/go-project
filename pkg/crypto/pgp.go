package crypto

import (
	"crypto/rand"
	"crypto/rsa"
)

var PGPPrivateKey *rsa.PrivateKey

func InitPGP() error {
	// Генерация ключей при первом запуске
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	PGPPrivateKey = privateKey
	return nil
}

func EncryptPGP(data string) (string, error) {
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, &PGPPrivateKey.PublicKey, []byte(data))
	if err != nil {
		return "", err
	}
	return string(encrypted), nil
}

func DecryptPGP(data string) (string, error) {
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, PGPPrivateKey, []byte(data))
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

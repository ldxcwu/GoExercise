package secret

import (
	"errors"
	"secret/encrypt"
)

type Vault struct {
	encodingKey string
	kvMap       map[string]string
}

func NewVault(encodingKey string) Vault {
	return Vault{
		encodingKey: encodingKey,
		kvMap:       make(map[string]string),
	}
}

func (v *Vault) Set(key, value string) error {
	encryptedVal, err := encrypt.Encrypt(v.encodingKey, value)
	if err != nil {
		panic(err)
	}
	v.kvMap[key] = encryptedVal
	return nil
}

func (v *Vault) Get(key string) (string, error) {
	encryptedVal, ok := v.kvMap[key]
	if !ok {
		return "", errors.New("no such value for that key")
	}
	plaintext, err := encrypt.Decrypt(v.encodingKey, encryptedVal)
	if err != nil {
		return "", err
	}
	return plaintext, nil
}

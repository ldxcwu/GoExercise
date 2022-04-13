package secret

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"secret/cipher"
	"secret/encrypt"
	"strings"
	"sync"
)

type Vault struct {
	encodingKey string
	filepath    string
	kvMap       map[string]string
	mux         sync.Mutex
}

func NewVault(encodingKey, filepath string) *Vault {
	return &Vault{
		encodingKey: encodingKey,
		filepath:    filepath,
	}
}

func (v *Vault) writeKV(w io.Writer) error {
	return json.NewEncoder(w).Encode(v.kvMap)
}

func (v *Vault) readKV(r io.Reader) error {
	return json.NewDecoder(r).Decode(&v.kvMap)
}

func (v *Vault) store() error {
	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	w, err := cipher.EncryptStreamWriter(v.encodingKey, f)
	if err != nil {
		return err
	}
	return v.writeKV(w)
}

func (v *Vault) load() error {
	f, err := os.Open(v.filepath)
	if err != nil {
		v.kvMap = make(map[string]string)
		return nil
	}
	defer f.Close()
	r, err := cipher.DecryptStreamReader(v.encodingKey, f)
	if err != nil {
		return err
	}
	return v.readKV(r)
}

func (v *Vault) loadKVMap() error {
	f, err := os.Open(v.filepath)
	if err != nil {
		v.kvMap = make(map[string]string)
		return nil
	}
	defer f.Close()
	// file is not nil
	// 1. load content from file(with encrypted json)
	// 2. decrypted file content to json
	// 3. decode json to v.kvMap
	var sb strings.Builder
	_, err = io.Copy(&sb, f)
	if err != nil {
		return err
	}
	decryptedJSON, err := encrypt.Decrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}
	dec := json.NewDecoder(strings.NewReader(decryptedJSON))
	err = dec.Decode(&v.kvMap)
	if err != nil {
		return err
	}
	return nil
}

func (v *Vault) storeKVMap() error {
	var sb strings.Builder
	err := json.NewEncoder(&sb).Encode(v.kvMap)
	if err != nil {
		return err
	}
	encryptedJSON, err := encrypt.Encrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}
	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(encryptedJSON)
	// _, err = fmt.Fprint(f, encryptedJSON)
	if err != nil {
		return err
	}
	return nil
}

func (v *Vault) Set(key, value string) error {
	v.mux.Lock()
	defer v.mux.Unlock()
	encryptedVal, err := cipher.Encrypt(v.encodingKey, value)
	if err != nil {
		panic(err)
	}
	err = v.load()
	if err != nil {
		return err
	}
	v.kvMap[key] = encryptedVal
	err = v.store()
	return err
}

func (v *Vault) Get(key string) (string, error) {
	v.mux.Lock()
	defer v.mux.Unlock()
	err := v.load()
	if err != nil {
		return "", err
	}
	encryptedVal, ok := v.kvMap[key]
	if !ok {
		return "", errors.New("no such value for that key")
	}
	plaintext, err := cipher.Decrypt(v.encodingKey, encryptedVal)
	if err != nil {
		return "", err
	}
	return plaintext, nil
}

func (v *Vault) Del(key string) error {
	v.mux.Lock()
	defer v.mux.Unlock()
	err := v.load()
	if err != nil {
		return err
	}
	_, ok := v.kvMap[key]
	if !ok {
		return nil
	}
	delete(v.kvMap, key)
	return v.store()
}

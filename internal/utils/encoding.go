package utils

import (
	"bytes"
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/minio/sio"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Encode(w io.Writer, source interface{}, encrypted bool, encryptKey [32]byte) error {
	if !encrypted {
		return json.NewEncoder(w).Encode(source)
	}

	buff := new(bytes.Buffer)
	err := json.NewEncoder(buff).Encode(source)
	if err != nil {
		return err
	}

	_, err = sio.Encrypt(w, buff, sio.Config{Key: encryptKey[:]})
	return err
}

func Decode(r io.Reader, target interface{}, encrypted bool, encryptKey [32]byte) error {
	if !encrypted {
		return json.NewDecoder(r).Decode(target)
	}

	// Have to decrypt at first
	buff := new(bytes.Buffer)
	_, err := sio.Decrypt(buff, r, sio.Config{Key: encryptKey[:]})
	if err != nil {
		return err
	}

	return json.NewDecoder(buff).Decode(target)

}

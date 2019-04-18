package main

import (
	"crypto/sha256"
	"log"
	"os"
	"testing"
)

const (
	testFolder      = "testdata"
	passphrase      = "test"
	jsonFilesConfig = "./configs/files.json"
	outputFolder    = "./decrypted-files"
)

func TestMain(m *testing.M) {
	err := os.Chdir(testFolder)
	if err != nil {
		log.Fatalln(err)
	}

	code := m.Run()

	// Remove outputFolder
	os.RemoveAll(outputFolder)

	os.Exit(code)
}

func TestDecryptor(t *testing.T) {
	app := App{
		config: config{
			PassPhrase:    passphrase,
			FilesJSONFile: jsonFilesConfig,
			OutputFolder:  outputFolder,
		},
		decodeKey: sha256.Sum256([]byte(passphrase)),
	}

	err := app.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	err = app.Decrypt()
	if err != nil {
		t.Fatal(err)
	}
}

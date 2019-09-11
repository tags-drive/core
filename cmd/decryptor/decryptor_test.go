package decryptor

import (
	"crypto/sha256"
	"os"
	"testing"

	clog "github.com/ShoshinNikita/log/v2"
)

const (
	testFolder      = "testdata"
	passphrase      = "test"
	jsonFilesConfig = "./var/files.json"
	dataFolder      = "./var/data"
)

var outputFolder = os.TempDir() + "/decrypted-files"

func TestMain(m *testing.M) {
	err := os.Chdir(testFolder)
	if err != nil {
		clog.Fatalln(err)
	}

	code := m.Run()

	// Remove outputFolder
	os.RemoveAll(outputFolder)

	os.Exit(code)
}

func TestDecryptor(t *testing.T) {
	app := app{
		config: config{
			PassPhrase:    passphrase,
			FilesJSONFile: jsonFilesConfig,
			DataFolder:    dataFolder,
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

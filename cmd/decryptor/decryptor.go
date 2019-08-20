package decryptor

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"os"
	"strconv"
	"sync"

	"github.com/ShoshinNikita/log/v2"
	"github.com/jessevdk/go-flags"
	"github.com/minio/sio"
	"github.com/pkg/errors"

	"github.com/tags-drive/core/internal/storage/files"
)

var log *clog.Logger

func init() {
	log = clog.NewProdLogger()
}

const (
	workersCount = 5
)

type config struct {
	PassPhrase string `long:"phrase" required:"true"`
	//
	FilesJSONFile string `long:"config-file" default:"./var/files.json"`
	DataFolder    string `long:"data-folder" default:"./var/data"`
	//
	OutputFolder string `short:"o" long:"output-folder" default:"./decrypted-files"`
	// We don't need DataFolder field because there are valid paths to encrypted files in FilesJSONFile
	// DataFolder string `long:"data-folder" default:"./data"`
}

type app struct {
	config config

	decodeKey [32]byte
}

func newApp() (*app, error) {
	app := new(app)
	_, err := flags.Parse(&app.config)
	if err != nil {
		return nil, err
	}

	app.decodeKey = sha256.Sum256([]byte(app.config.PassPhrase))

	return app, nil
}

// Prepare creates OutputFolder and checks FilesJSONFile
func (a *app) Prepare() error {
	f, err := os.Open(a.config.FilesJSONFile)
	if err != nil {
		return errors.Wrap(err, "invalid path to config file")
	}
	f.Close()

	err = os.MkdirAll(a.config.OutputFolder, 0666)
	return errors.Wrap(err, "can't create output folder")
}

func (a *app) Decrypt() error {
	filesList, err := a.getFilesList()
	if err != nil {
		return errors.Wrap(err, "invalid json file")
	}

	filesChan := make(chan files.File, 20)

	// Fill filesChan
	go func() {
		for i := range filesList {
			filesChan <- filesList[i]
		}
		close(filesChan)
	}()

	var wg sync.WaitGroup
	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var err error
			var input, output string
			for file := range filesChan {
				input = a.config.DataFolder + "/" + strconv.Itoa(file.ID)
				output = a.config.OutputFolder + "/" + file.Filename
				err = a.decryptAndSaveFile(input, output)
				if err != nil {
					log.Errorf("can't decrypt file %s: %s\n", file.Filename, err)
				}
			}
		}()
	}

	wg.Wait()

	return nil
}

func (a *app) getFilesList() (res []files.File, err error) {
	path := a.config.FilesJSONFile
	f, err := os.Open(path)
	if err != nil {
		return res, err
	}
	defer f.Close()

	// Decrypt file
	decryptedFile := new(bytes.Buffer)

	_, err = sio.Decrypt(decryptedFile, f, sio.Config{
		Key: a.decodeKey[:],
	})
	if err != nil {
		return res, errors.Wrap(err, "can't decrypt file")
	}

	// Decode json file
	filesObj := make(map[int]files.File)

	err = json.NewDecoder(decryptedFile).Decode(&filesObj)
	if err != nil {
		return res, err
	}

	// Convert to array
	res = make([]files.File, 0, len(filesObj))
	for _, f := range filesObj {
		res = append(res, f)
	}

	return res, nil
}

func (a *app) decryptAndSaveFile(encryptedFilePath, decryptedFilePath string) error {
	encryptedFile, err := os.Open(encryptedFilePath)
	if err != nil {
		return err
	}
	defer encryptedFile.Close()

	decryptedFile, err := os.Create(decryptedFilePath)
	if err != nil {
		return err
	}
	defer decryptedFile.Close()

	_, err = sio.Decrypt(decryptedFile, encryptedFile, sio.Config{
		Key: a.decodeKey[:],
	})

	return errors.Wrap(err, "can't decrypt file")
}

func StartDecryptor() <-chan struct{} {
	log = clog.NewProdConfig().PrintTime(false).Build()

	log.Infoln("init Decryptor")

	app, err := newApp()
	if err != nil {
		log.Fatalln(err)
	}

	err = app.Prepare()
	if err != nil {
		log.Fatalln(err)
	}

	log.Infoln("start decryption")

	err = app.Decrypt()
	if err != nil {
		log.Fatalln(err)
	}

	log.Infoln("decryption is finished")

	done := make(chan struct{})
	close(done)
	return done
}

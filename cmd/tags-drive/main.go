package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	clog "github.com/ShoshinNikita/log/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"github.com/tags-drive/core/internal/storage/files"
	"github.com/tags-drive/core/internal/storage/tags"
	"github.com/tags-drive/core/internal/web"
)

const version = "v0.9.0"

type config struct {
	Version string `ignored:"true"`

	Debug bool `envconfig:"DBG" default:"false"`

	// Web

	Port           string        `envconfig:"PORT" default:":80"`
	IsTLS          bool          `envconfig:"TLS" default:"false"`
	Login          string        `envconfig:"LOGIN" default:"user"`
	Password       string        `envconfig:"PSWRD" default:"qwerty"`
	SkipLogin      bool          `envconfig:"SKIP_LOGIN" default:"false"`     // Debug only
	MaxTokenLife   time.Duration `envconfig:"MAX_TOKEN_LIFE" default:"1440h"` // default is 60 days
	AuthCookieName string        `default:"auth"`                             // name of cookie that contains token

	// Storage

	Encrypt          bool     `envconfig:"ENCRYPT" default:"false"`
	PassPhraseString string   `envconfig:"PASS_PHRASE"`
	PassPhrase       [32]byte `ignored:"true"` // sha256 sum of PassPhraseString field

	StorageType string `envconfig:"STORAGE_TYPE" default:"json"`

	DataFolder          string        `default:"./var/data"`
	ResizedImagesFolder string        `default:"./var/data/resized"`
	TimeBeforeDeleting  time.Duration `envconfig:"TIME_BEFORE_DELETING" default:"168h"` // default is 168h = 7 days

	FilesJSONFile      string `default:"./var/files.json"`        // for files
	TagsJSONFile       string `default:"./var/tags.json"`         // for tags
	AuthTokensJSONFile string `default:"./var/auth_tokens.json"`  // for auth tokens
	ShareTokenJSONFile string `default:"./var/share_tokens.json"` // for share tokens
}

type App struct {
	config config

	server      *web.Server
	fileStorage *files.FileStorage
	tagStorage  *tags.TagStorage

	logger *clog.Logger
}

// PrepareNewApp parses globalConfig and inits services
func PrepareNewApp() (*App, error) {
	defer func() {
		// Reset sensitive env vars
		os.Setenv("LOGIN", "CLEARED")
		os.Setenv("PSWRD", "CLEARED")
		os.Setenv("PASS_PHRASE", "CLEARED")
	}()

	var cnf config
	err := envconfig.Process("", &cnf)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse Config")
	}

	// Checks
	if len(cnf.Port) > 0 && cnf.Port[0] != ':' {
		cnf.Port = ":" + cnf.Port
	}

	if cnf.Encrypt && cnf.PassPhraseString == "" {
		return nil, errors.New("wrong env config: PASS_PHRASE can't be empty with ENCRYPT=true")
	}

	if cnf.SkipLogin && !cnf.Debug {
		return nil, errors.New("wrong env config: SkipLogin can't be true in Production mode")
	}

	// Finish a config creation

	cnf.Version = version
	// Encrypt password
	cnf.Password = encryptPassword(cnf.Password)
	// Get PassPhrase
	cnf.PassPhrase = sha256.Sum256([]byte(cnf.PassPhraseString))
	cnf.PassPhraseString = ""

	app := &App{config: cnf}

	err = app.initServices()
	if err != nil {
		return nil, errors.Wrap(err, "can't init services")
	}

	return app, nil
}

const encryptRepeats = 11

func encryptPassword(s string) string {
	hash := sha256.Sum256([]byte(s))

	for i := 1; i < encryptRepeats; i++ {
		hash = sha256.Sum256([]byte(hex.EncodeToString(hash[:])))
	}

	return hex.EncodeToString(hash[:])
}

// initServices inits storages and server
func (app *App) initServices() error {
	app.logger = clog.NewProdLogger()
	if app.config.Debug {
		app.logger = clog.NewDevLogger()
	}

	var err error

	// File storage
	fileStorageConfig := files.Config{
		Debug:               app.config.Debug,
		DataFolder:          app.config.DataFolder,
		ResizedImagesFolder: app.config.ResizedImagesFolder,
		StorageType:         app.config.StorageType,
		FilesJSONFile:       app.config.FilesJSONFile,
		Encrypt:             app.config.Encrypt,
		PassPhrase:          app.config.PassPhrase,
		TimeBeforeDeleting:  app.config.TimeBeforeDeleting,
	}
	app.fileStorage, err = files.NewFileStorage(fileStorageConfig, app.logger)
	if err != nil {
		return errors.Wrap(err, "can't create new FileStorage")
	}

	// Tag storage
	tagStorageConfig := tags.Config{
		Debug:        app.config.Debug,
		StorageType:  app.config.StorageType,
		TagsJSONFile: app.config.TagsJSONFile,
		Encrypt:      app.config.Encrypt,
		PassPhrase:   app.config.PassPhrase,
	}
	app.tagStorage, err = tags.NewTagStorage(tagStorageConfig, app.logger)
	if err != nil {
		return errors.Wrap(err, "can't create new TagStorage")
	}

	// Web server
	serverConfig := web.Config{
		Debug:               app.config.Debug,
		DataFolder:          app.config.DataFolder,
		Port:                app.config.Port,
		IsTLS:               app.config.IsTLS,
		Login:               app.config.Login,
		Password:            app.config.Password,
		SkipLogin:           app.config.SkipLogin,
		AuthCookieName:      app.config.AuthCookieName,
		MaxTokenLife:        app.config.MaxTokenLife,
		AuthTokensJSONFile:  app.config.AuthTokensJSONFile,
		ShareTokensJSONFile: app.config.ShareTokenJSONFile,
		Encrypt:             app.config.Encrypt,
		PassPhrase:          app.config.PassPhrase,
		Version:             app.config.Version,
	}
	app.server, err = web.NewWebServer(serverConfig, app.fileStorage, app.tagStorage, app.logger)
	if err != nil {
		return errors.Wrap(err, "can't init WebServer")
	}

	return nil
}

func (app *App) Start() error {
	app.printConfig()

	app.logger.Infoln("start")

	shutdowned := make(chan struct{})

	// fatalErr is used when server went down
	fatalServerErr := make(chan struct{})

	// Goroutine to shutdown services
	go func() {
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		select {
		case <-term:
			app.logger.Warnln("interrupt signal")
		case <-fatalServerErr:
			// Nothing
		}

		// Shutdowns. Server must be first
		app.logger.Debugln("shutdown WebServer")
		err := app.server.Shutdown()
		if err != nil {
			app.logger.Warnf("can't shutdown server gracefully: %s\n", err)
		}

		app.logger.Debugln("shutdown FileStorage")
		err = app.fileStorage.Shutdown()
		if err != nil {
			app.logger.Warnf("can't shutdown FileStorage gracefully: %s\n", err)
		}

		app.logger.Debugln("shutdown TagStorage")
		err = app.tagStorage.Shutdown()
		if err != nil {
			app.logger.Warnf("can't shutdown TagStorage gracefully: %s\n", err)
		}

		close(shutdowned)
	}()

	app.fileStorage.StartBackgroundServices()

	if err := app.server.Start(); err != nil {
		app.logger.Errorf("server error: %s\n", err)
		close(fatalServerErr)
	}

	<-shutdowned

	app.logger.Infoln("stop")

	return nil
}

func (app *App) printConfig() {
	s := "Config:\n"

	vars := []struct {
		name string
		v    interface{}
	}{
		{"Debug", app.config.Debug},
		//
		{"Port", app.config.Port},
		{"TLS", app.config.IsTLS},
		{"Login", app.config.Login},
		{"SkipLogin", app.config.SkipLogin},
		//
		{"StorageType", app.config.StorageType},
		{"Encrypt", app.config.Encrypt},
	}

	for _, v := range vars {
		s += fmt.Sprintf("  * %-11s %v\n", v.name, v.v)
	}

	app.logger.WriteString(s)
}

func main() {
	log.SetFlags(0)
	log.Printf("Tags Drive %s - https://github.com/tags-drive\n", version)

	app, err := PrepareNewApp()
	if err != nil {
		log.Fatalln(err)
	}

	if err := app.Start(); err != nil {
		log.Fatalln(err)
	}
}

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

	auth "github.com/tags-drive/core/internal/storage/auth_tokens"
	"github.com/tags-drive/core/internal/storage/files"
	share "github.com/tags-drive/core/internal/storage/share_tokens"
	"github.com/tags-drive/core/internal/storage/tags"
	"github.com/tags-drive/core/internal/web"
)

const version = "v0.9.0"

type config struct {
	Version string `ignored:"true"`

	Debug bool `envconfig:"DEBUG" default:"false"`

	Web struct {
		Port     string `envconfig:"WEB_PORT" default:":80"`
		IsTLS    bool   `envconfig:"WEB_TLS" default:"false"`
		Login    string `envconfig:"WEB_LOGIN" default:"user"`
		Password string `envconfig:"WEB_PASSWORD" default:"qwerty"`
		// SkipLogin is availabel only in debug mode
		SkipLogin bool `envconfig:"WEB_SKIP_LOGIN" default:"false"`
		// The default value is 1440h (60 days)
		MaxTokenLife time.Duration `envconfig:"WEB_MAX_TOKEN_LIFE" default:"1440h"`
	}

	Storage struct {
		// TODO
		StorageType string `envconfig:"STORAGE_TYPE" default:"json"`

		Encrypt          bool   `envconfig:"STORAGE_ENCRYPT" default:"false"`
		PassPhraseString string `envconfig:"STORAGE_PASS_PHRASE"`
		// PassPhrase is a sha256 of PassPhraseString
		PassPhrase [32]byte

		TimeBeforeDeleting time.Duration `envconfig:"STORAGE_TIME_BEFORE_DELETING" default:"168h"` // default is 168h = 7 days
	}
}

// We use const vars for paths because the app is run in Docker container
const (
	DataFolder          = "./var/data"
	ResizedImagesFolder = "./var/data/resized"
	//
	FilesJSONFile       = "./var/files.json"        // for files
	TagsJSONFile        = "./var/tags.json"         // for tags
	AuthTokensJSONFile  = "./var/auth_tokens.json"  // for auth tokens
	ShareTokensJSONFile = "./var/share_tokens.json" // for share tokens
)

// AuthCookieName is a name of cookie that contains token
const AuthCookieName = "auth"

type App struct {
	config config

	fileStorage  *files.FileStorage
	tagStorage   *tags.TagStorage
	authService  *auth.AuthService
	shareService *share.ShareService
	server       *web.Server

	logger *clog.Logger
}

// PrepareNewApp parses globalConfig and inits services
func PrepareNewApp() (*App, error) {
	defer func() {
		// Reset sensitive env vars
		os.Setenv("WEB_LOGIN", "CLEARED")
		os.Setenv("WEB_PASSWORD", "CLEARED")
		os.Setenv("STORAGE_PASS_PHRASE", "CLEARED")
	}()

	var cnf config
	err := envconfig.Process("", &cnf)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse Config")
	}

	// Checks
	if len(cnf.Web.Port) > 0 && cnf.Web.Port[0] != ':' {
		cnf.Web.Port = ":" + cnf.Web.Port
	}

	if cnf.Storage.Encrypt && cnf.Storage.PassPhraseString == "" {
		return nil, errors.New("wrong env config: PASS_PHRASE can't be empty with ENCRYPT=true")
	}

	if cnf.Web.SkipLogin && !cnf.Debug {
		return nil, errors.New("wrong env config: SkipLogin can't be true in Production mode")
	}

	// Finish a config creation

	cnf.Version = version
	// Encrypt password
	cnf.Web.Password = encryptPassword(cnf.Web.Password)
	// Get PassPhrase4
	cnf.Storage.PassPhrase = sha256.Sum256([]byte(cnf.Storage.PassPhraseString))
	cnf.Storage.PassPhraseString = ""

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
		DataFolder:          DataFolder,
		ResizedImagesFolder: ResizedImagesFolder,
		StorageType:         app.config.Storage.StorageType,
		FilesJSONFile:       FilesJSONFile,
		Encrypt:             app.config.Storage.Encrypt,
		PassPhrase:          app.config.Storage.PassPhrase,
		TimeBeforeDeleting:  app.config.Storage.TimeBeforeDeleting,
	}
	app.fileStorage, err = files.NewFileStorage(fileStorageConfig, app.logger)
	if err != nil {
		return errors.Wrap(err, "can't create new FileStorage")
	}

	// Tag storage
	tagStorageConfig := tags.Config{
		Debug:        app.config.Debug,
		StorageType:  app.config.Storage.StorageType,
		TagsJSONFile: TagsJSONFile,
		Encrypt:      app.config.Storage.Encrypt,
		PassPhrase:   app.config.Storage.PassPhrase,
	}
	app.tagStorage, err = tags.NewTagStorage(tagStorageConfig, app.logger)
	if err != nil {
		return errors.Wrap(err, "can't create new TagStorage")
	}

	// Auth service
	authConfig := auth.Config{
		Debug:          app.config.Debug,
		TokensJSONFile: AuthTokensJSONFile,
		Encrypt:        app.config.Storage.Encrypt,
		PassPhrase:     app.config.Storage.PassPhrase,
		MaxTokenLife:   app.config.Web.MaxTokenLife,
	}
	app.authService, err = auth.NewAuthService(authConfig, app.logger)
	if err != nil {
		return errors.Wrap(err, "can't init new Auth Service")
	}

	// Share service
	shareConfig := share.Config{
		ShareTokenJSONFile: ShareTokensJSONFile,
		Encrypt:            app.config.Storage.Encrypt,
		PassPhrase:         app.config.Storage.PassPhrase,
	}
	app.shareService, err = share.NewShareStorage(shareConfig, app.fileStorage, app.logger)
	if err != nil {
		return errors.Wrap(err, "can't init new Share Service")
	}

	// Web server
	serverConfig := web.Config{
		Debug:          app.config.Debug,
		Port:           app.config.Web.Port,
		IsTLS:          app.config.Web.IsTLS,
		Login:          app.config.Web.Login,
		Password:       app.config.Web.Password,
		SkipLogin:      app.config.Web.SkipLogin,
		AuthCookieName: AuthCookieName,
		MaxTokenLife:   app.config.Web.MaxTokenLife,
		Version:        app.config.Version,
	}

	app.server, err = web.NewWebServer(serverConfig,
		app.fileStorage,
		app.tagStorage,
		app.authService,
		app.shareService,
		app.logger)
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

		// Shutdowns. Server must be the first

		app.logger.Debugln("shutdown WebServer")
		err := app.server.Shutdown()
		if err != nil {
			app.logger.Warnf("can't shutdown server gracefully: %s\n", err)
		}

		app.logger.Debugln("shutdown Auth Service")
		err = app.authService.Shutdown()
		if err != nil {
			app.logger.Warnf("can't shutdown Auth Service gracefully: %s\n", err)
		}

		app.logger.Debugln("shutdown Share Service")
		err = app.shareService.Shutdown()
		if err != nil {
			app.logger.Warnf("can't shutdown Share Service gracefully: %s\n", err)
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

	app.fileStorage.StartBackgroundJobs()
	app.authService.StartBackgroundJobs()

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
		{"Port", app.config.Web.Port},
		{"TLS", app.config.Web.IsTLS},
		{"Login", app.config.Web.Login},
		{"SkipLogin", app.config.Web.SkipLogin},
		//
		{"StorageType", app.config.Storage.StorageType},
		{"Encrypt", app.config.Storage.Encrypt},
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

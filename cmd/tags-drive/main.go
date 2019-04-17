package main

import (
	"crypto/sha256"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	clog "github.com/ShoshinNikita/log/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"github.com/tags-drive/core/internal/storage/files"
	"github.com/tags-drive/core/internal/storage/tags"
	"github.com/tags-drive/core/internal/web"
)

const version = "v0.6.0"

type globalConfig struct {
	Version string `ignored:"true"`

	// Debug defines is debug mode
	Debug bool `envconfig:"DBG" default:"false"`

	// Port of the server
	Port string `envconfig:"PORT" default:":80"`
	// IsTLS defines should the program use https
	IsTLS bool `envconfig:"TLS" default:"false"`
	// Login is a user login
	Login string `envconfig:"LOGIN" default:"user"`
	// Password is a user password
	Password string `envconfig:"PSWRD" default:"qwerty"`
	// SkipLogin let use Tags Drive without loginning (for Debug only)
	SkipLogin bool `envconfig:"SKIP_LOGIN" default:"false"`
	// MaxTokenLife defines the max lifetime of a token (2 months)
	MaxTokenLife time.Duration `envconfig:"MAX_TOKEN_LIFE" default:"1440h"`
	// AuthCookieName defines name of cookie that contains token
	AuthCookieName string `default:"auth"`

	// Encrypt defines, should the program encrypt files. False by default
	Encrypt bool `envconfig:"ENCRYPT" default:"false"`
	// PassPhrase is used to encrypt files. Key is a sha256 sum of env "PASS_PHRASE"
	PassPhrase [32]byte `ignored:"true"`

	// StorageType is a type of storage
	StorageType string `envconfig:"STORAGE_TYPE" default:"json"`

	// DataFolder is a folder where all files are kept
	DataFolder string `default:"./data"`
	// ResizedImagesFolder is a folder where all resized images are kept
	ResizedImagesFolder string `default:"./data/resized"`

	// FilesJSONFile is a json file with files information
	FilesJSONFile string `default:"./configs/files.json"`
	// TagsJSONFile is a json file with list of tags (with name and color)
	TagsJSONFile string `default:"./configs/tags.json"`
	// TokensJSONFile is a json file with list of tokens
	TokensJSONFile string `default:"./configs/tokens.json"`
}

type App struct {
	config globalConfig

	server      web.ServerInterface
	fileStorage files.FileStorageInterface
	tagStorage  tags.TagStorageInterface

	logger *clog.Logger
}

func NewApp() (*App, error) {
	var cnf globalConfig
	err := envconfig.Process("", &cnf)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse Config")
	}

	cnf.Version = version

	// Add PassPhrase
	if cnf.Encrypt {
		phrase := os.Getenv("PASS_PHRASE")
		if phrase == "" {
			return nil, errors.New("wrong env config: PASS_PHRASE can't be empty with ENCRYPT=true")
		}
		cnf.PassPhrase = sha256.Sum256([]byte(phrase))
	}

	app := &App{config: cnf}

	return app, nil
}

func (app *App) Start() error {
	lg := clog.NewProdLogger()
	if app.config.Debug {
		lg = clog.NewDevLogger()
	}

	app.logger = lg

	app.logger.Printf("Tags Drive %s (https://github.com/tags-drive)\n\n", app.config.Version)
	app.logger.Infoln("start")
	app.PrintConfig()

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
		Debug:          app.config.Debug,
		DataFolder:     app.config.DataFolder,
		Port:           app.config.Port,
		IsTLS:          app.config.IsTLS,
		Login:          app.config.Login,
		Password:       app.config.Password,
		SkipLogin:      app.config.SkipLogin,
		AuthCookieName: app.config.AuthCookieName,
		MaxTokenLife:   app.config.MaxTokenLife,
		TokensJSONFile: app.config.TokensJSONFile,
		Encrypt:        app.config.Encrypt,
		PassPhrase:     app.config.PassPhrase,
		Version:        app.config.Version,
	}
	app.server, err = web.NewWebServer(serverConfig, app.fileStorage, app.tagStorage, app.logger)
	if err != nil {
		return errors.Wrap(err, "can't init WebServer")
	}

	shutdowned := make(chan struct{})

	// fatalErr is used when server went down
	fatalServerErr := make(chan struct{})

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
		app.logger.Infoln("shutdown WebServer")
		err := app.server.Shutdown()
		if err != nil {
			app.logger.Warnf("can't shutdown server gracefully: %s\n", err)
		}

		app.logger.Infoln("shutdown FileStorage")
		err = app.fileStorage.Shutdown()
		if err != nil {
			app.logger.Warnf("can't shutdown FileStorage gracefully: %s\n", err)
		}

		app.logger.Infoln("shutdown TagStorage")
		err = app.tagStorage.Shutdown()
		if err != nil {
			app.logger.Warnf("can't shutdown TagStorage gracefully: %s\n", err)
		}

		close(shutdowned)
	}()

	err = app.server.Start()
	if err != nil {
		app.logger.Errorf("server error: %s\n", err)
		close(fatalServerErr)
	}

	<-shutdowned

	app.logger.Infoln("stop")

	return nil
}

func (app *App) PrintConfig() {
	app.logger.Infoln("Config:")

	vars := []struct {
		name string
		v    interface{}
	}{
		{"Debug", app.config.Debug},
		//
		{"Port", app.config.Port},
		{"TLS", app.config.IsTLS},
		{"Login", app.config.Login},
		{"Password", strings.Repeat("*", len(app.config.Password))},
		{"SkipLogin", app.config.SkipLogin},
		//
		{"StorageType", app.config.StorageType},
		{"Encrypt", app.config.Encrypt},
	}

	for _, v := range vars {
		app.logger.Printf("      * %-11s %v\n", v.name, v.v)
	}
}

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatalln(err)
	}

	if err := app.Start(); err != nil {
		log.Fatalln(err)
	}
}

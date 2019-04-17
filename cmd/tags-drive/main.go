package main

import (
	"crypto/sha256"
	"fmt"
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

type commonOptions struct {
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
	server      web.ServerInterface
	fileStorage files.FileStorageInterface
	tagStorage  tags.TagStorageInterface

	logger *clog.Logger

	options commonOptions
}

func NewApp() (*App, error) {
	var opts commonOptions
	err := envconfig.Process("", &opts)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse Config")
	}

	opts.Version = version

	// Add PassPhrase
	if opts.Encrypt {
		phrase := os.Getenv("PASS_PHRASE")
		if phrase == "" {
			return nil, errors.New("wrong env config: PASS_PHRASE can't be empty with ENCRYPT=true")
		}
		opts.PassPhrase = sha256.Sum256([]byte(phrase))
	}

	app := &App{options: opts}

	return app, nil
}

func (app *App) Start() error {
	lg := clog.NewProdLogger()
	if app.options.Debug {
		lg = clog.NewDevLogger()
	}

	lg.Printf("Tags Drive %s (https://github.com/tags-drive)\n\n", app.options.Version)

	lg.Infoln("start")

	// Print options
	lg.Infoln("options:")
	lg.Println(app.paramsToString())

	var err error

	app.logger = lg

	app.fileStorage, err = files.NewFileStorage(lg)
	if err != nil {
		return errors.Wrap(err, "can't create new FileStorage")
	}

	app.tagStorage, err = tags.NewTagStorage(lg)
	if err != nil {
		return errors.Wrap(err, "can't create new TagStorage")
	}

	app.server, err = web.NewWebServer(app.fileStorage, app.tagStorage, lg)
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
			lg.Warnln("interrupt signal")
		case <-fatalServerErr:
			// Nothing
		}

		// Shutdowns. Server must be first
		lg.Infoln("shutdown WebServer")
		err := app.server.Shutdown()
		if err != nil {
			lg.Warnf("can't shutdown server gracefully: %s\n", err)
		}

		lg.Infoln("shutdown FileStorage")
		err = app.fileStorage.Shutdown()
		if err != nil {
			lg.Warnf("can't shutdown FileStorage gracefully: %s\n", err)
		}

		lg.Infoln("shutdown TagStorage")
		err = app.tagStorage.Shutdown()
		if err != nil {
			lg.Warnf("can't shutdown TagStorage gracefully: %s\n", err)
		}

		close(shutdowned)
	}()

	err = app.server.Start()
	if err != nil {
		lg.Errorf("server error: %s\n", err)
		close(fatalServerErr)
	}

	<-shutdowned

	lg.Infoln("stop")

	return nil
}

func (app *App) paramsToString() string {
	s := "\n"

	vars := []struct {
		name string
		v    interface{}
	}{
		{"Port", app.options.Port},
		{"Login", app.options.Login},
		{"Password", strings.Repeat("*", len(app.options.Password))},
		{"TLS", app.options.IsTLS},
		{"Encrypt", app.options.Encrypt},
		{"StorageType", app.options.StorageType},
		{"Debug", app.options.Debug},
		{"SkipLogin", app.options.SkipLogin},
	}

	for _, v := range vars {
		s += fmt.Sprintf("\t* %-11s %v\n", v.name, v.v)
	}

	// Remove the last '\n'
	return s[:len(s)-1]
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

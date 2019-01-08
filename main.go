package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ShoshinNikita/log"

	"github.com/tags-drive/core/cmd"
	"github.com/tags-drive/core/internal/params"
	"github.com/tags-drive/core/internal/storage/files"
	"github.com/tags-drive/core/internal/storage/tags"
	"github.com/tags-drive/core/internal/web"
)

type App struct {
	Server      cmd.Server
	FileStorage cmd.FileStorageInterface
	TagStorage  cmd.TagStorageInterface
	Logger      *log.Logger
}

func main() {
	lg := log.NewLogger()
	lg.PrintColor(true)
	lg.PrintTime(true)

	if err := params.Parse(); err != nil {
		lg.Fatalln(err)
	}

	if params.Debug {
		lg.PrintErrorLine(true)
	}

	lg.Infoln("start")

	// Print params
	lg.Infoln("params:")
	lg.Println(paramsToString())

	var err error

	app := new(App)
	app.Logger = lg

	app.FileStorage, err = files.NewFileStorage(lg)
	if err != nil {
		lg.Fatalf("can't create new FileStorage: %s\n", err)
	}

	app.TagStorage, err = tags.NewTagStorage(lg)
	if err != nil {
		lg.Fatalf("can't create new TagStorage: %s\n", err)
	}

	app.Server, err = web.NewWebServer(app.FileStorage, app.TagStorage, lg)
	if err != nil {
		lg.Fatalf("can't init WebServer: %s", err)
	}

	shutdowned := make(chan struct{})

	go func() {
		term := make(chan os.Signal)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-term

		lg.Warnln("interrupt signal")

		// Shutdowns. Server must be first
		lg.Infoln("shutdown WebServer")
		err := app.Server.Shutdown()
		if err != nil {
			log.Warnf("can't shutdown server gracefully: %s\n", err)
		}

		lg.Infoln("shutdown FileStorage")
		err = app.FileStorage.Shutdown()
		if err != nil {
			log.Warnf("can't shutdown FileStorage gracefully: %s\n", err)
		}

		lg.Infoln("shutdown TagStorage")
		err = app.TagStorage.Shutdown()
		if err != nil {
			log.Warnf("can't shutdown TagStorage gracefully: %s\n", err)
		}

		close(shutdowned)
	}()

	err = app.Server.Start()
	if err != nil {
		lg.Errorf("server error: %s\n", err)
	}

	<-shutdowned

	lg.Infoln("stop")
}

func paramsToString() (s string) {
	vars := []struct {
		name string
		v    interface{}
	}{
		{"Port", params.Port},
		{"Login", params.Login},
		{"Password", strings.Repeat("*", len(params.Password))},
		{"TLS", params.IsTLS},
		{"Encrypt", params.Encrypt},
		{"StorageType", params.StorageType},
		{"Debug", params.Debug},
		{"SkipLogin", params.SkipLogin},
	}

	for _, v := range vars {
		s += fmt.Sprintf("\t* %-11s %v\n", v.name, v.v)
	}

	// Remove the last '\n'
	return s[:len(s)-1]
}

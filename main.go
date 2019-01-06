package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ShoshinNikita/log"

	"github.com/tags-drive/core/cmd"
	"github.com/tags-drive/core/internal/params"
	"github.com/tags-drive/core/internal/storage/files"
	"github.com/tags-drive/core/internal/storage/tags"
	"github.com/tags-drive/core/internal/web"
)

func paramsToString() (s string) {
	vars := []struct {
		name string
		v    interface{}
	}{
		{"Port", params.Port},
		{"Login", params.Login},
		{"Password", "******"},
		{"TLS", params.IsTLS},
		{"Encrypt", params.Encrypt},
		{"StorageType", params.StorageType},
		{"Debug", params.Debug},
		{"SkipLogin", params.SkipLogin},
	}

	for _, v := range vars {
		// "[INFO] " == 7 chars
		s += fmt.Sprintf("       * %s: %v\n", v.name, v.v)
	}

	// Remove the last '\n'
	return s[:len(s)-1]
}

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

	lg.Infoln("Start")

	// Print params
	lg.Infoln("Params:")
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

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		term := make(chan os.Signal)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-term

		lg.Warnln("Interrupt signal")

		cancel()
	}()

	err = app.Server.Start(ctx)
	if err != nil {
		lg.Fatalln(err)
	}

	lg.Infoln("Stop")
}

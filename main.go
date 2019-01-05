package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ShoshinNikita/log"

	"github.com/tags-drive/core/internal/params"
	"github.com/tags-drive/core/internal/storage"
	"github.com/tags-drive/core/internal/web"
	"github.com/tags-drive/core/internal/web/auth"
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

func main() {
	log.PrintColor(true)
	log.PrintTime(true)

	if params.Debug {
		log.PrintErrorLine(true)
	}

	log.Infoln("Start")

	// Print params
	log.Infoln("Params:")
	log.Println(paramsToString())

	err := storage.Init()
	if err != nil {
		log.Fatalln(err)
	}

	err = auth.Init()
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		term := make(chan os.Signal)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-term

		log.Warnln("Interrupt signal")

		cancel()
	}()

	err = web.Start(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	log.Infoln("Stop")
}

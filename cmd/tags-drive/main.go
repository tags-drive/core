package main

import (
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

	stopChan := make(chan struct{})
	errChan := make(chan error, 1)
	termChan := make(chan os.Signal, 1)

	go web.Start(stopChan, errChan)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-errChan:
		// Server is down
		log.Fatalln(err)
	case <-termChan:
		// We got SIGTERM, SIGKILL or SIGINT
		close(stopChan)
	}

	if err := <-errChan; err != nil {
		log.Fatalln(err)
	}

	log.Infoln("Stop")
}

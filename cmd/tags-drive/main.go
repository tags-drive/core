package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ShoshinNikita/log"

	"github.com/ShoshinNikita/tags-drive/internal/web"
)

func main() {
	log.ShowTime(false)
	log.PrintColor(true)

	log.Infoln("Start")

	stopChan := make(chan struct{})
	errChan := make(chan error, 1)
	termChan := make(chan os.Signal, 1)

	go web.Start(stopChan, errChan)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	select {
	case err := <-errChan:
		log.Println(err)
		close(stopChan)
	case <-termChan:
		close(stopChan)
	}

	if err := <-errChan; err != http.ErrServerClosed {
		log.Errorln(err)
	}

	log.Infoln("Stop")
}

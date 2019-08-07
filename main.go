package main

import (
	"github.com/tags-drive/core/cmd/app"
)

type Command func() <-chan struct{}

func main() {
	// The default command is app.StartApp
	var cmd Command = app.StartApp

	// TODO

	// Start
	<-cmd()
}

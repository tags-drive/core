package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tags-drive/core/cmd/app"
	"github.com/tags-drive/core/cmd/decryptor"
	"github.com/tags-drive/core/cmd/migrator"
)

type Command func() <-chan struct{}

func main() {
	commandList := map[string]Command{
		"":        app.StartApp, // the default command is app.StartApp
		"start":   app.StartApp,
		"decrypt": decryptor.StartDecryptor,
		"migrate": migrator.StartMigrator,
	}

	var (
		cmd           Command // command for execution
		passedCommand string
		ok            bool
	)

	args := os.Args[1:]
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		// Update the command
		passedCommand = args[0]
	}

	cmd, ok = commandList[passedCommand]
	if !ok {
		fmt.Println("[ERR] wrong command was passed")
		return
	}

	// Run
	<-cmd()
}

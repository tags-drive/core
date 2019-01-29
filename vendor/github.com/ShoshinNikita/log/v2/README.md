# colored-log

Pretty colored logger for Golang

## Preview

[![example](example/example.png)](example/example.go)

## Install

`go get github.com/ShoshinNikita/log`

## Example

[Example program](example/example.go)

```go
package main

import (
    clog "github.com/ShoshinNikita/log/v2"
)

func main() {
    // For prod use log.NewProdConfig() or log.NewProdLogger()
    // For dev use log.NewDevConfig() or log.NewDevLogger()
    c := &clog.Config{}

    l := c.PrintColor(true).PrintErrorLine(true).PrintTime(true).Build()

    l.Infoln("some info message")
    l.Warnln("some warn message")
    l.Errorln("some error message")

    l.WriteString("\n")

    l = c.PrintColor(true).PrintErrorLine(true).PrintTime(false).Build()
    l.Infoln("some info message")
    l.Warnln("some warn message")
    l.Errorln("some error message")

    l.WriteString("\n")

    l = c.PrintColor(true).PrintErrorLine(false).PrintTime(false).Build()
    l.Infoln("some info message")
    l.Warnln("some warn message")
    l.Errorln("some error message")
}
```

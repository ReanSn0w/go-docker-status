package main

import (
	"os"
	"time"

	"github.com/ReanSn0w/go-docker-status/internal/server"
	"github.com/ReanSn0w/gokit/pkg/app"
)

var (
	revision = "debug"
	opts     = struct {
		app.Debug

		Title string `short:"t" long:"title" env:"TITLE" default:"Docker" description:"Title of the status page"`
		Port  int    `short:"p" long:"port" env:"PORT" default:"8080" description:"Port to listen on"`
	}{}
)

func main() {
	app := app.New("Docker status page", revision, &opts)

	srv, err := server.New(app.Log(), opts.Title)
	if err != nil {
		app.Log().Logf("[ERROR] server init err: %s", err.Error())
		os.Exit(2)
	}

	app.Add(srv.Stop)
	srv.Start(app.CancelCause(), opts.Port)

	app.GracefulShutdown(time.Second * 3)
}

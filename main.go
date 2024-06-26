package main

import (
	"os"
	"time"

	"github.com/fortix/galera-ips/cmd"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// CGO_ENABLED=0 go build -ldflags="-s -w" -tags=netgo -installsuffix netgo -trimpath .

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC822})
	cmd.Execute()
}

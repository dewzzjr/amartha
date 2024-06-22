package main

import (
	"os"

	"github.com/dewzzjr/amartha/reconcile"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	reconcile.Reconcile(os.Args)
}

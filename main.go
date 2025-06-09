package main

import (
	"log/slog"

	"github.com/ToshihiroOgino/elib/log"
)

func main() {
	log.Init()

	slog.Info("Starting code generation")
}

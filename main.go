package main

import (
	"log/slog"

	"github.com/ToshihiroOgino/elib/log"
	"github.com/ToshihiroOgino/elib/usecase"
)

func main() {
	log.Init()

	slog.Info("Starting code generation")

	usecase.Create("aaaa@aaa.com", "password1234")
}

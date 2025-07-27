package usecase

import (
	"log/slog"

	"github.com/google/uuid"
)

func newUUID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		slog.Error("failed to generate UUID", "error", err)
		panic(err)
	}
	return id.String()
}

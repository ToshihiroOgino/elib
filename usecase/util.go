package usecase

import (
	"log/slog"
	"time"

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

func now() *time.Time {
	currentTime := time.Now()
	return &currentTime
}

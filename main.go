package main

import (
	"fmt"
	"log/slog"

	"github.com/ToshihiroOgino/elib/controller"
	"github.com/ToshihiroOgino/elib/env"
	"github.com/ToshihiroOgino/elib/infra/sqlite"
	"github.com/ToshihiroOgino/elib/log"
	"github.com/ToshihiroOgino/elib/secure"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Init()

	defer sqlite.CloseDB()

	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.Use(secure.SecurityMiddleware(), secure.CSRFMiddleware())
	router.SetHTMLTemplate(secure.LoadSecureTemplates())
	router.Static("/static", "./static")

	controller.NewController(router)

	serverAddr := fmt.Sprintf(":%d", env.Get().Port)
	slog.Info("Server starting", "address", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		slog.Error("Server failed to start", "error", err)
	}
}

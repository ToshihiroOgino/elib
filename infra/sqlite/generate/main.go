package main

import (
	"log/slog"

	"github.com/ToshihiroOgino/elib/infra/sqlite"
	"github.com/ToshihiroOgino/elib/log"
	"gorm.io/gen"
)

func main() {
	log.Init()
	slog.Info("Starting code generation")

	g := gen.NewGenerator(gen.Config{
		OutPath:       "generated/repository",
		ModelPkgPath:  "generated/domain",
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})
	db := sqlite.GetDB()
	g.UseDB(db)
	all := g.GenerateAllTable()
	g.ApplyBasic(all...)
	g.Execute()
}

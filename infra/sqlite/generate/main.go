package main

import (
	"log/slog"

	"github.com/ToshihiroOgino/elib/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	log.Init()
	slog.Info("Starting code generation")

	g := gen.NewGenerator(gen.Config{
		OutPath:       "repository",
		ModelPkgPath:  "domain",
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})
	db, err := gorm.Open(sqlite.Open("sqlite/db.sqlite3"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	g.UseDB(db)
	all := g.GenerateAllTable()
	g.ApplyBasic(all...)
	g.Execute()
}

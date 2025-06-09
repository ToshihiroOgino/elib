package main

import (
	"log/slog"

	"github.com/ToshihiroOgino/elib/log"
)

func main() {
	log.Init()

	slog.Info("Starting code generation")

	// g := gen.NewGenerator(gen.Config{
	// 	OutPath:       "./generated",
	// 	Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
	// 	WithUnitTest:  true,
	// 	FieldNullable: true,
	// })
	// db, err := gorm.Open(sqlite.Open("sqlite/db.sqlite3"), &gorm.Config{})
	// if err != nil {
	// 	panic(err)
	// }
	// g.UseDB(db)
	// g.ApplyBasic()

	// allUsers, _ := user.userRepository.FindAll()
	// 	for _, user := range allUsers {
	// 		println(user.Email) // Assuming User has an Email field
	// 	}
}

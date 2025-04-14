package migration

import (
	"embed"
	"fmt"

	"todo-list/internal/config"
	"todo-list/internal/pkg/store/db"
	"todo-list/pkg/migrator"
)

//go:embed files/*.sql
var migrationsFS embed.FS

func Run(db *db.DB, cfg *config.Config) error {
	migrator, err := migrator.NewMigrator(migrationsFS, cfg.DB)

	if err != nil {
		return err
	}

	err = migrator.Migrate(db.Connection())

	if err != nil {
		return err
	}

	fmt.Println("Migrations applied successfully!")

	return nil
}

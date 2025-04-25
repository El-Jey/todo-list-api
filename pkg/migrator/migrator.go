package migrator

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"todo-list/internal/config"
)

type Migrator struct {
	srcDriver source.Driver
	config    config.DB
}

func NewMigrator(sqlFiles embed.FS, config config.DB) (*Migrator, error) {
	d, err := iofs.New(sqlFiles, "files")

	if err != nil {
		return nil, err
	}

	return &Migrator{
		srcDriver: d,
		config:    config,
	}, nil
}

func (m *Migrator) Migrate(db *pgxpool.Pool) (err error) {
	driver, err := postgres.WithInstance(stdlib.OpenDBFromPool(db), &postgres.Config{
		MigrationsTable: m.config.MigrationsTable,
	})

	if err != nil {
		return fmt.Errorf("unable to create db instance: %v", err)
	}

	migrator, err := migrate.NewWithInstance("migration_embeded_sql_files",
		m.srcDriver,
		m.config.Database,
		driver)

	if err != nil {
		return fmt.Errorf("unable to create migration: %v", err)
	}

	defer func() {
		srcErr, dbErr := migrator.Close()

		if srcErr != nil && err == nil {
			err = fmt.Errorf("source error closing migrator: %v", srcErr)
		}

		if dbErr != nil && err == nil {
			err = fmt.Errorf("database error closing migrator: %v", dbErr)
		}
	}()

	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("unable to apply migrations %v", err)
	}

	return nil
}

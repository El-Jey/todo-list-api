package db

import (
	"context"

	"todo-list/internal/config"
	"todo-list/internal/pkg/store/db"
)

func InitDB(config config.DB, ctx context.Context) (*db.DB, error) {
	dbInstance, err := db.NewDB(config, ctx)

	if err != nil {
		return nil, err
	}

	return dbInstance, nil
}

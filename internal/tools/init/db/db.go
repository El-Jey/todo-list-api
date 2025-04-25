package db

import (
	"context"

	"todo-list/internal/config"
	"todo-list/internal/pkg/store/db"
)

func InitDB(ctx context.Context, config config.DB) (*db.DB, error) {
	dbInstance, err := db.NewDB(ctx, config)

	if err != nil {
		return nil, err
	}

	return dbInstance, nil
}

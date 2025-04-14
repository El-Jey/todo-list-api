package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	migration "todo-list/db/migrations"
	apiServer "todo-list/internal/app/api"
	"todo-list/internal/config"
	dbInit "todo-list/internal/tools/init/db"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)

	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Get()

	if err != nil {
		return err
	}

	dbInstance, err := dbInit.InitDB(cfg.DB, ctx)

	if err != nil {
		return err
	}

	defer dbInstance.Close()

	if err := migration.Run(dbInstance, &cfg); err != nil {
		return err
	}

	srv := apiServer.NewServer(&cfg, dbInstance, ctx)

	go func() {
		if err := srv.Listen(fmt.Sprintf("%s:%s", cfg.Server.API.Host, cfg.Server.API.Port)); err != nil {
			fmt.Fprintf(os.Stderr, "error listening and serving Api server: %s\n", err)
		}
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()

		if err := srv.ShutdownWithContext(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()

	wg.Wait()

	return nil
}

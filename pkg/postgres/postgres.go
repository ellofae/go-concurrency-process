package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"

	"github.com/ellofae/go-concurrency-process/config"
	"github.com/ellofae/go-concurrency-process/internal/utils"
	"github.com/ellofae/go-concurrency-process/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func OpenPoolConnection(ctx context.Context, cfg *config.Config) (conn *pgxpool.Pool) {
	logger := logger.GetLogger()

	err := utils.ConnectionAttemps(func() error {
		var err error

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		conn, err = pgxpool.New(ctx, fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.PostgresDB.User,
			cfg.PostgresDB.Password,
			cfg.PostgresDB.Host,
			cfg.PostgresDB.Port,
			cfg.PostgresDB.DBName,
			cfg.PostgresDB.SSLmode,
		))

		return err
	}, 3, time.Duration(2)*time.Second)

	if err != nil {
		logger.Error("Didn't manage to make connection with database", "message", err.Error())
		os.Exit(1)
	}

	logger.Info("Database connection is established successfully.")

	return conn
}

func RunMigrationsUp(ctx context.Context, cfg *config.Config) {
	logger := logger.GetLogger()

	db_conn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.PostgresDB.User,
		cfg.PostgresDB.Password,
		cfg.PostgresDB.Host,
		cfg.PostgresDB.Port,
		cfg.PostgresDB.DBName,
		cfg.PostgresDB.SSLmode,
	)

	migration, err := migrate.New("file://migrations", db_conn)
	if err != nil {
		logger.Error("Unable to get a migrate instance", "error", err.Error())
		os.Exit(1)
	}

	err = migration.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			logger.Warn("No changes while migrating")
			return
		}

		logger.Error("Unable to migrate up", "error", err.Error())
		os.Exit(1)
	}
	logger.Info("Migrations are up successfully.")
}

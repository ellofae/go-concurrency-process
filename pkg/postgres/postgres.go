package postgres

import (
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

type PostgresConn struct {
	conn *sqlx.DB
}

func OpenDatabaseConnection() (*PostgresConn, error) {
	db_conn, err := PostgresConnection()
	if err != nil {
		return nil, err
	}

	return &PostgresConn{conn: db_conn}, nil
}

func PostgresConnection() (*sqlx.DB, error) {
	godotenv.Load(".env")

	db, err := sqlx.Connect("pgx", os.Getenv("DATABASE_CONNECTION_URI"))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

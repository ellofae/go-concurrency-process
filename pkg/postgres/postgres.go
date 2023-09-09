package postgres

import (
	"os"

	"github.com/jmoiron/sqlx"
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
	db, err := sqlx.Connect("pgx", os.Getenv("DATABASE_CONNECTION_URI"))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

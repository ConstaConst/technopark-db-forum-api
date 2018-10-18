package postgres

import (
	"github.com/jackc/pgx"
	"io/ioutil"
)

func Connect() (*pgx.ConnPool, error) {
	var pgxConfig = pgx.ConnConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "forum_api",
		User:     "forum_api",
		Password: "q",
	}

	dbConn, err := pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig:     pgxConfig,
			MaxConnections: 8,
		})
	if err != nil {
		return nil, err
	}

	err = createTables(dbConn)
	if err != nil {
		return nil, err
	}

	return dbConn, err
}

func createTables(dbConn *pgx.ConnPool) error {
	initSql, err := ioutil.ReadFile("postgres/init.sql")
	if err != nil {
		return err
	}

	_, err = dbConn.Exec(string(initSql))
	if err != nil {
		return err
	}

	return nil
}

func Close(dbConn *pgx.ConnPool) {
	dbConn.Close()
}

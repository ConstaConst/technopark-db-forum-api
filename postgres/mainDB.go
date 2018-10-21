package postgres

import (
	"github.com/jackc/pgx"
	"io/ioutil"
	"log"
	"strconv"
)

type DBConn struct {
	pool *pgx.ConnPool
}

func MakeDBConn() (DBConn, error) {
	var pgxConfig = pgx.ConnConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "forum_api",
		User:     "forum_api",
		Password: "q",
	}

	pool, err := pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig:     pgxConfig,
			MaxConnections: 8,
		})
	if err != nil {
		return DBConn{}, err
	}

	log.Println("Conected: address=" + pgxConfig.Host + ":" +
		strconv.Itoa(int(pgxConfig.Port)) +
		" db=" + pgxConfig.Database + " user=" + pgxConfig.User)

	return DBConn{pool}, nil
}

func (conn *DBConn) InitDBTables() error {
	initSql, err := ioutil.ReadFile("postgres/init.sql")
	if err != nil {
		return err
	}

	_, err = conn.pool.Exec(string(initSql))
	if err != nil {
		return err
	}

	return nil
}

func (conn *DBConn) Close() {
	conn.pool.Close()

	log.Println("Connection to DB is closed")
}

package postgres

import (
	"github.com/ConstaConst/technopark-db-forum-api/models"
	"github.com/ConstaConst/technopark-db-forum-api/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"log"
)

func (conn *DBConn) ShowStatus(
	params operations.StatusParams) middleware.Responder {
	tx, _ := conn.pool.Begin()
	defer tx.Rollback()

	log.Println("Service status")

	row := tx.QueryRow(`SELECT usersNumber, forumsNumber, threadsNumber, postsNumber
                             FROM service
                             LIMIT 1`)

	var service models.Status
	row.Scan(&service.User, &service.Forum, &service.Thread, &service.Post)

	tx.Commit()

	return operations.NewStatusOK().WithPayload(&service)
}

func (conn *DBConn) UpdateStatus() {
	log.Println("Update status")

	_, err := conn.pool.Exec("UPDATE service SET (usersNumber, forumsNumber, threadsNumber, postsNumber)" +
		"=((SELECT count(*) FROM users), (SELECT count(*) FROM forums), (SELECT count(*) FROM threads)," +
		" (SELECT count(*) FROM posts))")
	checkError(err)
}

func (conn *DBConn) ClearService(
	params operations.ClearParams) middleware.Responder {
	tx, _ := conn.pool.Begin()
	defer tx.Rollback()

	log.Println("Service clear")

	_, err := tx.Exec("TRUNCATE TABLE forums, threads, users, posts, " +
		"votes, users_in_forums, service RESTART IDENTITY CASCADE")
	checkError(err)

	_, err = tx.Exec("INSERT INTO service VALUES (0, 0, 0, 0)")
	checkError(err)

	tx.Commit()

	return operations.NewClearOK()
}

package postgres

import (
	"fmt"
	"github.com/ConstaConst/technopark-db-forum-api/models"
	"github.com/ConstaConst/technopark-db-forum-api/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"log"
)

func (conn *DBConn) CreateThread(
	params operations.ThreadCreateParams) middleware.Responder {
	tx, _ := conn.pool.Begin()
	defer tx.Rollback()

	log.Printf("Create Thread. Params: in forum=%s, author=%s, slug=%s\n",
		params.Slug, params.Thread.Author, params.Thread.Slug)

	author, err := getUser(tx, params.Thread.Author)
	if err != nil {
		notFoundUserError := models.Error{Message: fmt.Sprintf(
			"Can't find author by nickname=%s", params.Thread.Author)}

		tx.Rollback()
		return operations.NewThreadCreateNotFound().WithPayload(
			&notFoundUserError)
	}

	forum, err := getForum(tx, params.Slug)
	if err != nil {
		notFoundForumError := models.Error{Message: fmt.Sprintf(
			"Can't find forum by slug=%s", params.Slug)}

		tx.Rollback()
		return operations.NewThreadCreateNotFound().WithPayload(
			&notFoundForumError)
	}

	var thread models.Thread
	var slug *string
	if params.Thread.Slug != "" {
		thread, err = getThread(tx, params.Thread.Slug)
		slug = &params.Thread.Slug
	} else {
		thread, err = getThreadByFAT(tx, params.Thread)
		slug = nil
	}
	if err == nil {
		log.Println("Thread slug=", thread.Slug, "already exists")

		tx.Rollback()
		return operations.NewThreadCreateConflict().WithPayload(&thread)
	}

	params.Thread.Author = author.Nickname
	params.Thread.Forum = forum.Slug

	var created interface{}
	if params.Thread.Created != nil {
		created = params.Thread.Created
	} else {
		created = "NOW()"
	}

	row := tx.QueryRow(`INSERT INTO threads
					VALUES (DEFAULT, $1, $2, $3, $4, $5, $6, $7)
					RETURNING id`,
		slug, params.Thread.Title, params.Thread.Message,
		params.Thread.Author, params.Thread.Forum, created,
		params.Thread.Votes)
	err = row.Scan(&params.Thread.ID)
	checkError(err)

	tx.Commit()

	log.Println("Thread id=", params.Thread.ID, " is created")

	return operations.NewThreadCreateCreated().WithPayload(params.Thread)
}

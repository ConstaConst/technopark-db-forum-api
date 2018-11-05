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

	_, err = tx.Exec("UPDATE forums SET threadsNumber=threadsNumber+1 WHERE slug=$1",
		forum.Slug)
	checkError(err)
	_, err = tx.Exec("UPDATE service SET threadsNumber=threadsNumber+1")
	checkError(err)

	_, err = tx.Exec("INSERT INTO users_in_forums (forum, nickname) VALUES ($1, $2)"+
		" ON CONFLICT (forum, nickname) DO NOTHING",
		forum.Slug, params.Thread.Author)

	tx.Commit()

	log.Println("Thread id=", params.Thread.ID, " is created")

	return operations.NewThreadCreateCreated().WithPayload(params.Thread)
}

func (conn *DBConn) GetOneThread(params operations.ThreadGetOneParams) middleware.Responder {
	tx, _ := conn.pool.Begin()
	defer tx.Rollback()

	log.Println("Get one thread=", params.SlugOrID)

	thread, err := getThread(tx, params.SlugOrID)
	if err != nil {
		notFoundThreadError := models.Error{Message: fmt.Sprintf(
			"Can't find thread =%s", params.SlugOrID)}

		tx.Rollback()
		return operations.NewThreadGetOneNotFound().WithPayload(
			&notFoundThreadError)
	}

	tx.Commit()

	return operations.NewThreadGetOneOK().WithPayload(&thread)
}

func (conn *DBConn) UpdateThread(
	params operations.ThreadUpdateParams) middleware.Responder {
	tx, _ := conn.pool.Begin()
	defer tx.Rollback()

	log.Println("Update Thread=", params.SlugOrID)

	thread, err := getThread(tx, params.SlugOrID)
	if err != nil {
		tx.Rollback()

		log.Println("Can't find thread by slug or id=", params.SlugOrID)

		notFoundThreadError := models.Error{Message: fmt.Sprintf(
			"Can't find thread by slug or id=%s", params.SlugOrID)}

		return operations.NewThreadUpdateNotFound().WithPayload(
			&notFoundThreadError)
	}

	if params.Thread.Title == "" && params.Thread.Message == "" {
		return operations.NewThreadUpdateOK().WithPayload(&thread)
	}

	var args []interface{}
	query := "UPDATE threads SET "

	if params.Thread.Title != "" {
		args = append(args, params.Thread.Title)
		query += fmt.Sprintf("title=$%d, ", len(args))
		thread.Title = params.Thread.Title
	}
	if params.Thread.Message != "" {
		args = append(args, params.Thread.Message)
		query += fmt.Sprintf("message=$%d ", len(args))
		thread.Message = params.Thread.Message
	}

	args = append(args, thread.ID)
	query += fmt.Sprintf("WHERE id=$%d", len(args))

	_, err = tx.Exec(query, args...)
	checkError(err)

	tx.Commit()

	log.Println("Thread=", params.SlugOrID, "is updated")

	return operations.NewThreadUpdateOK().WithPayload(&thread)
}

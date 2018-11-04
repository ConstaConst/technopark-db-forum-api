package postgres

import (
	"fmt"
	"github.com/ConstaConst/technopark-db-forum-api/models"
	"github.com/ConstaConst/technopark-db-forum-api/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx/pgtype"
	"log"
)

func (conn *DBConn) CreateForum(
	params operations.ForumCreateParams) middleware.Responder {
	tx, _ := conn.pool.Begin()
	defer tx.Rollback()

	log.Printf("Create Forum. Params: %v\n", params.Forum)

	user, err := getUser(tx, params.Forum.User)
	if err != nil {
		notFoundUserError := models.Error{Message: fmt.Sprintf(
			"Can't find user by nickname=%s", params.Forum.User)}

		tx.Rollback()
		return operations.NewForumCreateNotFound().WithPayload(
			&notFoundUserError)
	}

	forum, err := getForum(tx, params.Forum.Slug)
	if err == nil {
		log.Println("Forum slug=", params.Forum.Slug, "already exists")

		tx.Rollback()
		return operations.NewForumCreateConflict().WithPayload(&forum)
	}

	_, err = tx.Exec(`INSERT INTO forums 
					VALUES ($1, $2, $3, $4, $5)`,
		params.Forum.Slug, params.Forum.Title,
		user.Nickname, params.Forum.Posts, params.Forum.Threads)
	checkError(err)

	_, err = tx.Exec(
		"UPDATE service SET forumsNumber=forumsNumber+1")
	checkError(err)

	tx.Commit()

	params.Forum.User = user.Nickname

	log.Println("Forum is created")

	return operations.NewForumCreateCreated().WithPayload(params.Forum)
}

func (conn *DBConn) GetOneForum(
	params operations.ForumGetOneParams) middleware.Responder {
	tx, _ := conn.pool.Begin()
	defer tx.Rollback()

	log.Printf("Get One Forum. Params: %s\n", params.Slug)

	forum, err := getForum(tx, params.Slug)
	if err != nil {
		log.Println("Forum slug=", params.Slug, "isn't found")

		tx.Rollback()
		notFoundForumError := models.Error{Message: fmt.Sprintf(
			"Can't find forum by slag=%s", params.Slug)}
		return operations.NewForumGetOneNotFound().WithPayload(
			&notFoundForumError)
	}
	tx.Commit()

	log.Println("Forum slug=", params.Slug, "is found")

	return operations.NewForumGetOneOK().WithPayload(&forum)
}

func (conn *DBConn) GetForumThreads(
	params operations.ForumGetThreadsParams) middleware.Responder {
	tx, _ := conn.pool.Begin()
	defer tx.Rollback()

	log.Printf("Get Forum Threads. Params: %s\n", params.Slug)

	_, err := getForum(tx, params.Slug)
	if err != nil {
		log.Println("Forum slug=", params.Slug, "isn't found")

		tx.Rollback()
		notFoundForumError := models.Error{Message: fmt.Sprintf(
			"Can't find forum by slag=%s", params.Slug)}
		return operations.NewForumGetThreadsNotFound().WithPayload(
			&notFoundForumError)
	}

	var args []interface{}
	query := "SELECT id, slug, title, message, author, forum, created, " +
		"votesNumber " +
		"FROM threads " +
		"WHERE forum = $1 "
	args = append(args, params.Slug)
	if params.Since != nil {
		args = append(args, params.Since.String())
		if params.Desc != nil && *params.Desc {
			query += fmt.Sprintf("AND created <= $%d::timestamptz ", len(args))
		} else {
			query += fmt.Sprintf("AND created >= $%d::timestamptz ", len(args))
		}
	}
	query += "ORDER BY created "
	if params.Desc != nil && *params.Desc {
		query += "DESC "
	}
	if params.Limit != nil {
		args = append(args, *params.Limit)
		query += fmt.Sprintf("LIMIT $%d", len(args))
	}

	rows, err := tx.Query(query, args...)
	checkError(err)

	threads := models.Threads{}
	for rows.Next() {
		thread := models.Thread{}
		fetchedCreated := pgtype.Timestamptz{}
		fetchedSlug := pgtype.Text{}
		err = rows.Scan(&thread.ID, &fetchedSlug, &thread.Title, &thread.Message,
			&thread.Author, &thread.Forum, &fetchedCreated, &thread.Votes)
		checkError(err)
		t := strfmt.NewDateTime()
		err = t.Scan(fetchedCreated.Time)
		checkError(err)
		thread.Created = &t

		thread.Slug = fetchedSlug.String

		threads = append(threads, &thread)
	}

	log.Println("Threads are fetched:", len(threads))

	return operations.NewForumGetThreadsOK().WithPayload(threads)
}

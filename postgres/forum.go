package postgres

import (
	"fmt"
	"github.com/ConstaConst/technopark-db-forum-api/models"
	"github.com/ConstaConst/technopark-db-forum-api/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"log"
)

func (conn *DBConn) CreateForum(params operations.ForumCreateParams) middleware.Responder {
	tx, _ := conn.pool.Begin()
	defer tx.Rollback()

	log.Printf("Create Forum. Params: %v\n", params.Forum)

	user, err := getUser(tx, params.Forum.User)
	if err != nil {
		notFoundUserError := models.Error{Message: fmt.Sprintf(
			"Can't find user by nickname=%s", params.Forum.User)}

		tx.Commit()
		return operations.NewForumCreateNotFound().WithPayload(
			&notFoundUserError)
	}

	forum, err := getForum(tx, params.Forum.Slug)
	if err == nil {
		log.Println("Forum slug=", params.Forum.Slug, "already exists")

		tx.Commit()
		return operations.NewForumCreateConflict().WithPayload(&forum)
	}

	tx.Exec(`INSERT INTO forums 
					VALUES ($1, $2, $3, $4, $5)`,
		params.Forum.Slug, params.Forum.Title,
		user.Nickname, params.Forum.Posts, params.Forum.Threads)
	tx.Commit()

	params.Forum.User = user.Nickname

	log.Println("Forum is created")

	return operations.NewForumCreateCreated().WithPayload(params.Forum)
}


func (conn *DBConn) GetOneForum(params operations.ForumGetOneParams) middleware.Responder {
	tx, _ := conn.pool.Begin()
	defer tx.Rollback()

	log.Printf("Get One Forum. Params: %s\n", params.Slug)

	forum, err := getForum(tx, params.Slug)
	if err != nil {
		log.Println("Forum slug=", params.Slug, "isn't found")

		tx.Commit()
		notFoundForumError := models.Error{Message: fmt.Sprintf(
			"Can't find forum by slag=%s", params.Slug)}
		return operations.NewForumGetOneNotFound().WithPayload(
			&notFoundForumError)
	}
	tx.Commit()

	log.Println("Forum slug=", params.Slug, "is found")

	return operations.NewForumGetOneOK().WithPayload(&forum)
}
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

func (conn *DBConn) CreatePosts(
	params operations.PostsCreateParams) middleware.Responder {
	tx, _ := conn.pool.Begin()
	defer tx.Rollback()

	log.Printf("Create Posts. Params: in thread=%s, posts count=%d\n",
		params.SlugOrID, len(params.Posts))

	if len(params.Posts) == 0 {
		log.Println("Nothing to insert")

		return operations.NewPostsCreateCreated().WithPayload(models.Posts{})
	}

	thread, err := getThread(tx, params.SlugOrID)
	if err != nil {
		tx.Rollback()

		log.Println("Can't find thread by slug or id=", params.SlugOrID)

		notFoundThreadError := models.Error{Message: fmt.Sprintf(
			"Can't find thread by slug or id=%s", params.SlugOrID)}

		return operations.NewPostsCreateNotFound().WithPayload(
			&notFoundThreadError)
	}

	var args []interface{}
	query := "INSERT INTO posts (author, message, forum, thread, parent, path) " +
		"VALUES "
	j := 1
	for i, post := range params.Posts {
		if post.Parent != 0 {

			_, err = getPost(tx, post.Parent)
			if err != nil {
				tx.Rollback()

				log.Println("Can't find post parent:", post.Parent)

				notFoundParentError := models.Error{Message: fmt.Sprintf(
					"Can't find post parent=%d", post.Parent)}
				return operations.NewPostsCreateConflict().WithPayload(
					&notFoundParentError)
			}
		}

		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,(SELECT path FROM posts WHERE id=$%d)"+
			"||(SELECT currval('posts_id_seq')))", j, j+1, j+2, j+3, j+4, j+5)
		if i < len(params.Posts)-1 {
			query += ", "
		}
		j += 6
		args = append(args, post.Author, post.Message, thread.Forum, thread.ID, post.Parent, post.Parent)
	}
	query += "RETURNING *;"

	rows, err := tx.Query(query, args...)
	checkError(err)

	posts := models.Posts{}
	for rows.Next() {
		post := models.Post{}
		fetchedCreated := pgtype.Timestamptz{}
		var _tPath interface{}
		err = rows.Scan(&post.ID, &post.Author, &post.Message, &post.Forum,
			&post.Thread, &post.Parent, &fetchedCreated, &post.IsEdited, _tPath)
		checkError(err)
		t := strfmt.NewDateTime()
		err = t.Scan(fetchedCreated.Time)
		checkError(err)
		post.Created = &t

		posts = append(posts, &post)
	}

	tx.Commit()

	log.Println("Posts are created:", len(posts))

	return operations.NewPostsCreateCreated().WithPayload(posts)
}

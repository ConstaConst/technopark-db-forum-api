package postgres

import (
	"fmt"
	"github.com/ConstaConst/technopark-db-forum-api/models"
	"github.com/ConstaConst/technopark-db-forum-api/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
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

func (conn *DBConn) GetPosts(
	params operations.ThreadGetPostsParams) middleware.Responder {
	tx, _ := conn.pool.Begin()
	defer tx.Rollback()

	log.Printf("Get posts in thread=%s. Params: desc=%d, limit=%d, since=%d,"+
		" sort=%s.",
		params.SlugOrID, params.Desc, params.Limit,
		params.Since, *params.Sort)

	thread, err := getThread(tx, params.SlugOrID)
	if err != nil {
		tx.Rollback()

		log.Println("Can't find thread by slug or id=", params.SlugOrID)

		notFoundThreadError := models.Error{Message: fmt.Sprintf(
			"Can't find thread by slug or id=%s", params.SlugOrID)}

		return operations.NewThreadGetPostsNotFound().WithPayload(
			&notFoundThreadError)
	}

	var args []interface{}
	var query string
	if params.Sort != nil {
		switch *params.Sort {
		case "flat":
			query, args = makeQueryFlat(&params, thread.ID)
		case "tree":
			query, args = makeQueryTree(&params, thread.ID)
		case "parent_tree":
			query, args = makeQueryParentTree(tx, &params, thread.ID)
		default:
			return middleware.NotImplemented("Unkwnown sort type")
		}
	}
	if query == "" {
		log.Println("Not found available posts")

		return operations.NewThreadGetPostsOK().WithPayload(models.Posts{})
	}
	log.Println("Query: ", query)
	log.Println("Args:", args)
	rows, err := tx.Query(query, args...)
	checkError(err)

	posts := models.Posts{}
	for rows.Next() {
		post := models.Post{}
		fetchedCreated := pgtype.Timestamptz{}
		err = rows.Scan(&post.Author, &fetchedCreated, &post.Forum, &post.ID,
			&post.Message, &post.Thread, &post.Parent)
		checkError(err)
		t := strfmt.NewDateTime()
		err = t.Scan(fetchedCreated.Time)
		checkError(err)
		post.Created = &t

		posts = append(posts, &post)
	}
	tx.Commit()

	log.Println("Fetched posts ", len(posts))

	return operations.NewThreadGetPostsOK().WithPayload(posts)
}

func makeQueryFlat(params *operations.ThreadGetPostsParams, thread int32) (string, []interface{}) {
	log.Println("Make flat query")

	var args []interface{}
	query := "SELECT author, created, forum, id, message, thread, parent " +
		"FROM posts " +
		"WHERE thread=$1 "
	args = append(args, thread)

	if params.Since != nil {
		args = append(args, *params.Since)
		if params.Desc != nil && *params.Desc {
			query += fmt.Sprintf("AND id < $%d ", len(args))
		} else {
			query += fmt.Sprintf("AND id > $%d ", len(args))
		}
	}
	query += "ORDER BY id "
	if params.Desc != nil && *params.Desc {
		query += "DESC "
	}

	if params.Limit != nil {
		args = append(args, *params.Limit)
		query += fmt.Sprintf("LIMIT $%d", len(args))
	}

	return query, args
}

func makeQueryTree(params *operations.ThreadGetPostsParams, thread int32) (string, []interface{}) {
	log.Println("Make tree query")

	var args []interface{}
	query := "SELECT author, created, forum, id, message, thread, parent " +
		"FROM posts " +
		"WHERE thread=$1 "
	args = append(args, thread)

	if params.Since != nil {
		args = append(args, *params.Since)
		if params.Desc != nil && *params.Desc {
			query += fmt.Sprintf(
				"AND path < (select path from posts where id=$%d)", len(args))
		} else {
			query += fmt.Sprintf(
				"AND path > (select path from posts where id=$%d) ", len(args))
		}
	}
	query += "ORDER BY path "
	if params.Desc != nil && *params.Desc {
		query += "DESC "
	}

	if params.Limit != nil {
		args = append(args, *params.Limit)
		query += fmt.Sprintf("LIMIT $%d", len(args))
	}

	return query, args
}

func makeQueryParentTree(tx *pgx.Tx, params *operations.ThreadGetPostsParams, thread int32) (string, []interface{}) {
	log.Println("Make parent tree query")

	var args []interface{}
	query := "SELECT author, created, forum, id, message, thread, parent " +
		"FROM posts " +
		"WHERE thread=$1 "
	args = append(args, thread)

	if params.Since != nil || params.Desc != nil || params.Limit != nil {
		roots, _ := getRootPosts(tx, params, thread)
		if len(roots) == 0 {
			return "", args
		}
		rootsStr := ""
		for i, root := range roots {
			args = append(args, root)
			rootsStr += fmt.Sprintf("path && ARRAY[$%d]::bigint[] ", len(args))

			if i < len(roots)-1 {
				rootsStr += "OR "
			}
		}
		query += "AND (" + rootsStr + ")"
	}

	if params.Desc != nil && *params.Desc {
		query += "ORDER BY path[1] DESC, path "
	} else {
		query += "ORDER BY path "
	}

	return query, args
}

func getRootPosts(tx *pgx.Tx, params *operations.ThreadGetPostsParams, thread int32) ([]int64, error) {
	subQuery := "SELECT id " +
		"FROM posts " +
		"WHERE thread=$1 AND array_length(path, 1)=1 "
	var subArgs []interface{}
	subArgs = append(subArgs, thread)

	if params.Since != nil {
		subArgs = append(subArgs, *params.Since)
		if params.Desc != nil && *params.Desc {
			subQuery += fmt.Sprintf("AND path < ARRAY[(select path[1] from posts where id=$%d)] ", len(subArgs))
		} else {
			subQuery += fmt.Sprintf("AND path > ARRAY[(select path from posts where id=$%d)] ", len(subArgs))
		}
	}

	subQuery += "ORDER BY path "
	if params.Desc != nil && *params.Desc {
		subQuery += "DESC "
	}

	if params.Limit != nil {
		subArgs = append(subArgs, *params.Limit)
		subQuery += fmt.Sprintf("LIMIT $%d ", len(subArgs))
	}

	log.Println("subQuery:", subQuery)
	log.Println("subArgs:", subArgs)

	rows, err := tx.Query(subQuery, subArgs...)
	checkError(err)

	var roots []int64

	for rows.Next() {
		var root int64
		err = rows.Scan(&root)
		checkError(err)

		roots = append(roots, root)
	}

	return roots, nil
}

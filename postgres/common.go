package postgres

import (
	"github.com/ConstaConst/technopark-db-forum-api/models"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"log"
	"strconv"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func getUser(tx *pgx.Tx, nickname string) (models.User, error) {
	row := tx.QueryRow(`SELECT nickname, fullname, email, about 
								FROM users 
								WHERE nickname = $1`,
		nickname)
	user := models.User{}
	err := row.Scan(&user.Nickname, &user.Fullname, &user.Email, &user.About)
	if err != nil {
		log.Println("User=", nickname, "isn't found.", err)

		return models.User{}, err
	}

	log.Println("User=", user.Nickname, "is found")

	return user, nil
}

func getForum(tx *pgx.Tx, slug string) (models.Forum, error) {
	row := tx.QueryRow(`SELECT slug, title, "user", postsnumber, threadsnumber
								FROM forums
								WHERE slug = $1`,
		slug)
	forum := models.Forum{}

	err := row.Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Posts,
		&forum.Threads)
	if err != nil {
		log.Println("Forum=", slug, "isn't found.", err)

		return models.Forum{}, err
	}

	log.Println("Forum=", forum.Slug, "is found")

	return forum, nil
}

// How do make identification the thread without slug or id?
func getThreadByFAT(tx *pgx.Tx, reqThread *models.Thread) (models.Thread,
	error) {
	row := tx.QueryRow(`SELECT id, slug, title, message, author, forum,
								created, votesNumber
								FROM threads
								WHERE forum=$1 AND author=$2 AND title=$3`,
		reqThread.Forum, reqThread.Author, reqThread.Title)
	thread := models.Thread{}
	fetchedCreated := pgtype.Timestamptz{}
	fetchedSlug := pgtype.Text{}
	err := row.Scan(&thread.ID, &fetchedSlug, &thread.Title, &thread.Message,
		&thread.Author, &thread.Forum, &fetchedCreated, &thread.Votes)
	if err != nil {
		log.Printf("Thread with forum=%s, author=%s, title=%s isn't found: ",
			reqThread.Forum, reqThread.Author, reqThread.Title)
		log.Println(err)

		return models.Thread{}, err
	}
	t := strfmt.NewDateTime()
	err = t.Scan(fetchedCreated.Time)
	checkError(err)
	thread.Created = &t

	thread.Slug = fetchedSlug.String

	log.Println("Thread=", thread.ID, "is found")

	return thread, nil
}

func getThread(tx *pgx.Tx, slugOrId string) (models.Thread, error) {
	var queryType string
	if _, err := strconv.Atoi(slugOrId); err != nil {
		queryType = "slug=$1"
		log.Println("Fetch thread by slug:", slugOrId)
	} else {
		queryType = "id=$1::bigint"
		log.Println("Fetch thread by id:", slugOrId)
	}
	row := tx.QueryRow(`SELECT id, slug, title, message, author, forum,
								created, votesNumber
								FROM threads
								WHERE `+queryType,
		slugOrId)
	thread := models.Thread{}
	fetchedCreated := pgtype.Timestamptz{}
	fetchedSlug := pgtype.Text{}
	err := row.Scan(&thread.ID, &fetchedSlug, &thread.Title, &thread.Message,
		&thread.Author, &thread.Forum, &fetchedCreated, &thread.Votes)
	if err != nil {
		log.Println("Thread=", slugOrId, "isn't found.", err)

		return models.Thread{}, err
	}
	t := strfmt.NewDateTime()
	err = t.Scan(fetchedCreated.Time)
	checkError(err)
	thread.Created = &t

	thread.Slug = fetchedSlug.String

	log.Println("Thread=", thread.ID, "is found")

	return thread, nil
}

func getPost(tx *pgx.Tx, id int64) (models.Post, error) {
	row := tx.QueryRow(`SELECT id, author, message, forum, thread, parent, 
								created, isEdited
								FROM posts
								WHERE id = $1`,
		id)
	post := models.Post{}
	fetchedCreated := pgtype.Timestamptz{}
	err := row.Scan(&post.ID, &post.Author, &post.Message, &post.Forum,
		&post.Thread, &post.Parent, &fetchedCreated, &post.IsEdited)
	if err != nil {
		log.Println("Post=", id, "isn't found.", err)

		return models.Post{}, err
	}

	t := strfmt.NewDateTime()
	err = t.Scan(fetchedCreated.Time)
	checkError(err)
	post.Created = &t

	log.Println("Forum=", post.ID, "is found")

	return post, nil
}

func getVote(tx *pgx.Tx, nickname string, thread int32) (models.Vote, error) {
	row := tx.QueryRow(`SELECT nickname, voice
                             FROM votes
                             WHERE nickname=$1 AND thread=$2`,
		nickname, thread)

	vote := models.Vote{}
	err := row.Scan(&vote.Nickname, &vote.Voice)
	if err != nil {
		log.Println("Voice nickname=", nickname, "thread=", thread,
			"isn't found")

		return models.Vote{}, err
	}

	log.Println("Voice nickname=", nickname, "thread=", thread,
		"is found and =", vote.Voice)

	return vote, nil
}

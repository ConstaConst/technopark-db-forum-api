package postgres

import (
	"github.com/ConstaConst/technopark-db-forum-api/models"
	"github.com/jackc/pgx"
	"log"
	"strconv"
)

func checkError(err *error) {
	if err != nil {
		log.Fatalln(err)
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

func getForum(tx *pgx.Tx, slug string) (models.Forum, error)  {
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

func getThread(tx *pgx.Tx, slugOrId string) (models.Thread, error) {
	var queryType string
	if _, err := strconv.Atoi(slugOrId); err != nil {
		queryType = "slug=$1"
	} else {
		queryType = "id=$1::bigint"
	}
	row := tx.QueryRow(`SELECT id, slug, title, message, author, forum,
								created, votesNumber
								FROM threads
								WHERE ` + queryType,
								slugOrId)
	thread := models.Thread{}
	err := row.Scan(&thread.ID, &thread.Slug, &thread.Title, &thread.Message,
		&thread.Author, &thread.Forum, &thread.Created, &thread.Title)
	if err != nil {
		log.Println("Thread=", slugOrId, "isn't found.", err)

		return models.Thread{}, err
	}

	log.Println("Forum=", thread.ID, "is found")

	return thread, nil
}

func getPost(tx *pgx.Tx, id int64) (models.Post, error) {
	row := tx.QueryRow(`SELECT id, author, message, forum, thread, parent, 
								created, isEdited
								FROM posts
								WHERE id = $1`,
		id)
	post := models.Post{}
	err := row.Scan(&post.ID, &post.Author, &post.Message, &post.Forum,
		&post.Thread, &post.Parent, &post.Created, &post.IsEdited)
	if err != nil {
		log.Println("Post=", id, "isn't found.", err)

		return models.Post{}, err
	}

	log.Println("Forum=", post.ID, "is found")

	return post, nil
}

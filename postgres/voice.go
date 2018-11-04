package postgres

import (
	"fmt"
	"github.com/ConstaConst/technopark-db-forum-api/models"
	"github.com/ConstaConst/technopark-db-forum-api/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"log"
)

func (conn *DBConn) ThreadVote(
	params operations.ThreadVoteParams) middleware.Responder {
	tx, _ := conn.pool.Begin()
	defer tx.Rollback()

	log.Println("Make vote =", params.Vote.Voice,
		"in thread = ", params.SlugOrID)

	thread, err := getThread(tx, params.SlugOrID)
	if err != nil {
		tx.Rollback()

		log.Println("Can't find thread by slug or id=", params.SlugOrID)

		notFoundThreadError := models.Error{Message: fmt.Sprintf(
			"Can't find thread by slug or id=%s", params.SlugOrID)}

		return operations.NewThreadVoteNotFound().WithPayload(
			&notFoundThreadError)
	}

	oldVote, err := getVote(tx, params.Vote.Nickname, thread.ID)
	var deltaVoice int32
	if err != nil {
		deltaVoice = params.Vote.Voice

		_, err = tx.Exec(`INSERT INTO votes
                      VALUES ($1, $2, $3)`,
			params.Vote.Nickname, thread.ID, params.Vote.Voice)
	} else {
		deltaVoice = params.Vote.Voice - oldVote.Voice

		_, err = tx.Exec(`UPDATE votes
                              SET voice=$3
                              WHERE nickname=$1 AND thread=$2`,
			params.Vote.Nickname, thread.ID, params.Vote.Voice)
	}
	checkError(err)

	thread.Votes += deltaVoice
	_, err = tx.Exec(`UPDATE threads
                              SET votesnumber=$1
                              WHERE id=$2;`,
		thread.Votes, thread.ID)
	checkError(err)

	tx.Commit()

	log.Println("Votes updated:", thread.Votes)

	return operations.NewThreadVoteOK().WithPayload(&thread)
}

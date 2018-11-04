package postgres

import (
	"fmt"
	"github.com/ConstaConst/technopark-db-forum-api/models"
	"github.com/ConstaConst/technopark-db-forum-api/restapi/operations"
	"github.com/go-openapi/runtime/middleware"
	"log"
)

func (conn *DBConn) CreateUser(
	params operations.UserCreateParams) middleware.Responder {
	tx, _ := conn.pool.Begin()
	defer tx.Rollback()

	log.Println("Create User. Params: nickname=", params.Nickname,
		", email=", params.Profile.Email)

	rows, _ := tx.Query(`SELECT nickname, fullname, email, about 
								FROM users 
								WHERE nickname = $1 OR email = $2`,
		params.Nickname, params.Profile.Email)
	existedUsers := models.Users{}
	for rows.Next() {
		eUser := models.User{}
		rows.Scan(&eUser.Nickname, &eUser.Fullname, &eUser.Email, &eUser.About)
		existedUsers = append(existedUsers, &eUser)
	}

	if len(existedUsers) != 0 {
		log.Println("User exists")

		tx.Rollback()
		return operations.NewUserCreateConflict().WithPayload(existedUsers)
	}

	_, err := tx.Exec(`INSERT INTO users 
					VALUES ($1, $2, $3, $4)`,
		params.Nickname, params.Profile.Fullname,
		params.Profile.Email, params.Profile.About)
	checkError(err)

	_, err = tx.Exec(
		"UPDATE service SET usersNumber=usersNumber+1")
	checkError(err)

	tx.Commit()

	log.Println("User is created")

	params.Profile.Nickname = params.Nickname

	return operations.NewUserCreateCreated().WithPayload(params.Profile)
}

func (conn *DBConn) GetOneUser(
	params operations.UserGetOneParams) middleware.Responder {
	tx, _ := conn.pool.Begin()
	defer tx.Rollback()

	log.Println("Get One User. Params: nickname=", params.Nickname)

	row := tx.QueryRow(`SELECT nickname, fullname, email, about 
								FROM users 
								WHERE nickname = $1`,
		params.Nickname)
	user := models.User{}
	err := row.Scan(&user.Nickname, &user.Fullname, &user.Email, &user.About)
	if err != nil {
		log.Println("User=", params.Nickname, "is'n found")

		notFoundUserError := models.Error{Message: fmt.Sprintf(
			"Can't find user by nickname=%s", params.Nickname)}

		tx.Rollback()
		return operations.NewUserGetOneNotFound().WithPayload(
			&notFoundUserError)
	}
	tx.Commit()

	log.Println("User=", params.Nickname, "is found")

	return operations.NewUserGetOneOK().WithPayload(&user)
}

func (conn *DBConn) UpdateUser(
	params operations.UserUpdateParams) middleware.Responder {
	tx, _ := conn.pool.Begin()
	defer tx.Rollback()

	log.Println("Update User. Params: nickname=", params.Nickname)

	rows, _ := tx.Query(`SELECT nickname, fullname, email, about
								FROM users
								WHERE nickname=$1 or email=$2`,
		params.Nickname, params.Profile.Email)
	userCount := 0
	curUser := models.User{}
	for rows.Next() {
		userCount++
		if userCount <= 1 {
			rows.Scan(&curUser.Nickname, &curUser.Fullname, &curUser.Email,
				&curUser.About)
		}
	}
	if userCount > 1 {
		log.Println("Conflict. Found", userCount, "users with nickname=",
			params.Nickname, "or email=", params.Profile.Email)

		tx.Rollback()

		conflictUsersError := models.Error{Message: fmt.Sprintf(
			"Exist %d users with nickname=%s or email=%s",
			userCount, params.Nickname, params.Profile.Email)}
		return operations.NewUserUpdateConflict().WithPayload(
			&conflictUsersError)
	}
	if userCount == 0 {
		log.Println("User=", params.Nickname, "is'n found.")

		tx.Rollback()

		notFoundUserError := models.Error{Message: fmt.Sprintf(
			"Can't find user by nickname=%s", params.Nickname)}
		return operations.NewUserUpdateNotFound().WithPayload(
			&notFoundUserError)
	}
	if params.Profile.Email != "" || params.Profile.Fullname != "" ||
		params.Profile.About != "" {
		if params.Profile.Email != "" {
			curUser.Email = params.Profile.Email
		}
		if params.Profile.Fullname != "" {
			curUser.Fullname = params.Profile.Fullname
		}
		if params.Profile.About != "" {
			curUser.About = params.Profile.About
		}

		tx.Exec(`UPDATE users
						SET fullname=$1, email=$2, about=$3
						WHERE nickname=$4`,
			curUser.Fullname, curUser.Email, curUser.About,
			curUser.Nickname)
	}

	tx.Commit()

	return operations.NewUserUpdateOK().WithPayload(&curUser)
}

package users

import (
	"fmt"

	"github.com/swjang1214/bookstore_users-api/datasources/mysql/users_db"
	"github.com/swjang1214/bookstore_users-api/utils/date_utils"
	"github.com/swjang1214/bookstore_users-api/utils/errors"
)

const (
	queryInsertUser = "INSERT INTO users (first_name, last_name, email, date_created) VALUES (?,?,?,?)"
)

var (
	usersDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestError {

	if err := users_db.Client.Ping(); err != nil {
		panic(err)
	}

	result := usersDB[user.ID]
	if result == nil {
		return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.ID))
	}

	user.ID = result.ID
	user.FirstName = result.FirstName
	user.LastName = result.LastName
	user.Email = result.Email
	user.DateCreated = result.DateCreated

	return nil
}

func (user *User) Save() *errors.RestError {

	//! 실행 쿼리를 만들고
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		fmt.Println(err)
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	user.DateCreated = date_utils.GetNowString()

	fmt.Println(users_db.Client)
	//! 데이터를 넣어 실행
	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if err != nil {
		return errors.NewInternalServerError(
			fmt.Sprintf("error When trying to save user : %s", err.Error()))
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(
			fmt.Sprintf("error When trying to save user : %s", err.Error()))
	}
	user.ID = userId

	// current := usersDB[user.ID]
	// if current != nil {
	// 	if current.Email == user.Email {
	// 		return errors.NewBadRequestError(fmt.Sprintf("user %d already exists", user.ID))
	// 	}
	// 	return errors.NewBadRequestError(fmt.Sprintf("user %d already exists", user.ID))
	// }

	// user.DateCreated = date_utils.GetNowString()

	// usersDB[user.ID] = user
	return nil

}

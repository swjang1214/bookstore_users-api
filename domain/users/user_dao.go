package users

import (
	"fmt"

	"github.com/swjang1214/bookstore_users-api/datasources/mysql/users_db"
	"github.com/swjang1214/bookstore_users-api/utils/errors"
	"github.com/swjang1214/bookstore_users-api/utils/mysql_utils"
)

const (
	uniqueKeyword         = "users.email"
	errorNoRows           = "no rows in result set"
	queryInsertUser       = "INSERT INTO users (first_name, last_name, email, date_created, status, password) VALUES (?,?,?,?,?,?);"
	queryGetUser          = "SELECT id, first_name, last_name, email, date_created, status, password FROM users WHERE id=?;"
	queryUpdateUser       = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser       = "DELETE FROM users WHERE id=?;"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
)

var (
	usersDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestError {

	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		fmt.Println(err)
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.ID)
	if result == nil {
		return errors.NewInternalServerError("query error")
	}

	/*
		results, err := stmt.Query(user.ID)
		if err != nil {
			return errors.NewInternalServerError(err.Error())
		}
		defer results.Close()
	*/

	if err := result.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status, &user.Password); err != nil {
		// if strings.Contains(err.Error(), errorNoRows) {
		// 	return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.ID))
		// }

		// return errors.NewInternalServerError(
		// 	fmt.Sprintf("error when trying to get user %d : %s", user.ID, err.Error()))
		return mysql_utils.ParseError(err)
	}

	// result := usersDB[user.ID]
	// if result == nil {
	// 	return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.ID))
	// }

	// user.ID = result.ID
	// user.FirstName = result.FirstName
	// user.LastName = result.LastName
	// user.Email = result.Email
	// user.DateCreated = result.DateCreated

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

	fmt.Println(users_db.Client)
	//! 데이터를 넣어 실행
	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if err != nil {

		return mysql_utils.ParseError(err)
		// sqlErr, ok := err.(*mysql.MySQLError)
		// if !ok {
		// 	return errors.NewInternalServerError(
		// 		fmt.Sprintf("error When trying to save user : %s", err.Error()))
		// }
		// // fmt.Println(sqlErr.Number)
		// // fmt.Println(sqlErr.Message)
		// switch sqlErr.Number {
		// case 1062:
		// 	return errors.NewBadRequestError(fmt.Sprintf("email %s already exist", user.Email))

		// }

		// // if strings.Contains(err.Error(), uniqueKeyword) {
		// // 	return errors.NewBadRequestError(fmt.Sprintf("email %s already exist", user.Email))
		// // }
		// return errors.NewInternalServerError(
		// 	fmt.Sprintf("error When trying to save user : %s", err.Error()))
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		return mysql_utils.ParseError(err)
		// return errors.NewInternalServerError(
		// 	fmt.Sprintf("error When trying to save user : %s", err.Error()))
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

func (user *User) Update() *errors.RestError {

	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		fmt.Println(err)
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.ID)
	if err != nil {
		return mysql_utils.ParseError(err)
	}

	return nil
}

func (user *User) Delete() *errors.RestError {

	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		fmt.Println(err)
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID)
	if err != nil {
		return mysql_utils.ParseError(err)
	}

	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestError) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		fmt.Println(err)
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		fmt.Println(err)
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close()
	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			if err != nil {
				return nil, mysql_utils.ParseError(err)
			}
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil
}

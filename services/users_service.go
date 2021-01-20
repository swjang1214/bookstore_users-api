package services

import (
	"github.com/swjang1214/bookstore_users-api/domain/users"
	"github.com/swjang1214/bookstore_users-api/utils/date_utils"
	"github.com/swjang1214/bookstore_users-api/utils/errors"
)

func GetUser(userId int64) (*users.User, *errors.RestError) {
	user := users.User{
		ID: userId,
	}
	err := user.Get()
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user users.User) (*users.User, *errors.RestError) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	user.Status = users.StatusActive
	user.DateCreated = date_utils.GetNowDBString()
	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUser(isPartial bool, user users.User) (*users.User, *errors.RestError) {

	if isPartial == false {
		if err := user.Validate(); err != nil {
			return nil, err
		}
	}

	//! 먼저 임시 사용자 정보를 얻어 온다
	current := users.User{
		ID: user.ID,
	}
	err := current.Get()
	if err != nil {
		return nil, err
	}

	if isPartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}
		if user.LastName != "" {
			current.LastName = user.LastName
		}
		if user.Email != "" {
			current.Email = user.Email
		}
	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.Email = user.Email
		current.DateCreated = user.DateCreated
	}

	if err := current.Update(); err != nil {
		return nil, err
	}

	return &user, nil
}

func DeleteUser(userId int64) *errors.RestError {
	user := &users.User{ID: userId}
	return user.Delete()
}

func FindByStatus(status string) ([]users.User, *errors.RestError) {
	dao := &users.User{}
	return dao.FindByStatus(status)
}

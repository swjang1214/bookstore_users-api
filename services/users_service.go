package services

import (
	"fmt"

	"github.com/swjang1214/bookstore_users-api/domain/users"
	"github.com/swjang1214/bookstore_users-api/utils/crypto_utils"
	"github.com/swjang1214/bookstore_users-api/utils/date_utils"
	"github.com/swjang1214/bookstore_users-api/utils/errors"
)

type userService struct{}
type userServiceInterface interface {
	GetUser(int64) (*users.User, *errors.RestError)
	CreateUser(users.User) (*users.User, *errors.RestError)
	UpdateUser(bool, users.User) (*users.User, *errors.RestError)
	DeleteUser(int64) *errors.RestError
	SearchUser(string) (users.Users, *errors.RestError)
	LoginUser(users.LoginRequest) (*users.User, *errors.RestError)
}

var (
	UserService userServiceInterface = &userService{}
)

func (*userService) LoginUser(req users.LoginRequest) (*users.User, *errors.RestError) {
	dao := &users.User{
		Email:    req.Email,
		Password: req.Password,
	}

	dao.Password = crypto_utils.GetMd5(dao.Password)
	err := dao.GetByEmailAndPassword()
	if err != nil {
		return nil, err
	}

	return dao, nil

}

func (*userService) GetUser(userId int64) (*users.User, *errors.RestError) {
	user := users.User{
		ID: userId,
	}
	err := user.Get()
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (*userService) CreateUser(user users.User) (*users.User, *errors.RestError) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	user.Status = users.StatusActive
	user.DateCreated = date_utils.GetNowDBString()
	user.Password = crypto_utils.GetMd5(user.Password)
	fmt.Println(user.Password)
	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (*userService) UpdateUser(isPartial bool, user users.User) (*users.User, *errors.RestError) {

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

func (*userService) DeleteUser(userId int64) *errors.RestError {
	user := &users.User{ID: userId}
	return user.Delete()
}

func (*userService) SearchUser(status string) (users.Users, *errors.RestError) {
	dao := &users.User{}
	return dao.FindByStatus(status)
}

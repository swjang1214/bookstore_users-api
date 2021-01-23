package users

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/swjang1214/bookstore_users-api/utils/crypto_utils"
	"github.com/swjang1214/bookstore_users-api/utils/errors"
)

const (
	invitationHashLayout = "%s_%s_%v"
	StatusActive         = "active"
	StatusUsed           = "used"
)

type Invitation struct {
	ID          int64  `json:"id"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Hash        string `json:"hash"`
	Status      string `json:"status"`
}

func (i *Invitation) GetNewHash() string {
	return crypto_utils.GetMd5(fmt.Sprintf(invitationHashLayout, i.Email, i.DateCreated, rand.Uint64()))
}

type User struct {
	ID          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Password    string `json:"password"`
}

type Users []User

func Validate(user *User) *errors.RestError {
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return errors.NewBadRequestError("invalid email address")
	}
	return nil
}

func (user *User) Validate() *errors.RestError {
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return errors.NewBadRequestError("invalid email address")
	}
	user.Password = strings.TrimSpace(user.Password)
	if user.Email == "" {
		return errors.NewBadRequestError("invalid password")
	}
	return nil
}

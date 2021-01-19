package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/swjang1214/bookstore_users-api/domain/users"
	"github.com/swjang1214/bookstore_users-api/services"
	"github.com/swjang1214/bookstore_users-api/utils/errors"
)

var (
	counter int
)

func CreateUser(c *gin.Context) {
	var user users.User

	// 1. 2와 결과는 같다
	// bytes, err := ioutil.ReadAll(c.Request.Body)
	// if err != nil {
	// 	//TODO : Handle Error
	// 	return
	// }
	// if err := json.Unmarshal(bytes, &user); err != nil {
	// 	//TODO : Handle json Error
	// 	return
	// }

	// 2. 보완점 실제 json의 키들과의 값 테스트가 필요할 듯
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	result, restErr := services.CreateUser(user)
	if restErr != nil {
		c.JSON(restErr.Status, restErr)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func GetUser(c *gin.Context) {
	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestError("user id should be a number")
		c.JSON(err.Status, &err)
		return
	}
	user, getErr := services.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, user)
}

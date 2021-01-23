package users_controller

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

func TestServiceInterface() {
	//services
}

func getUserId(userIdParam string) (int64, *errors.RestError) {
	userId, userErr := strconv.ParseInt(userIdParam, 10, 64)
	if userErr != nil {
		return 0, errors.NewBadRequestError("user id should be a number")
	}
	return userId, nil
}

func Create(c *gin.Context) {
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

	result, restErr := services.UserService.CreateUser(user)
	if restErr != nil {
		c.JSON(restErr.Status, restErr)
		return
	}

	c.JSON(http.StatusCreated, result)
}

func Get(c *gin.Context) {

	userId, userErr := getUserId(c.Param("user_id"))
	if userErr != nil {
		c.JSON(userErr.Status, &userErr)
		return
	}

	// userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	// if userErr != nil {
	// 	err := errors.NewBadRequestError("user id should be a number")
	// 	c.JSON(err.Status, &err)
	// 	return
	// }
	user, getErr := services.UserService.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, user.Marshal(c.GetHeader("X-Public") == "true"))
}

func Update(c *gin.Context) {
	// userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	// if userErr != nil {
	// 	err := errors.NewBadRequestError("user id should be a number")
	// 	c.JSON(err.Status, &err)
	// 	return
	// }
	userId, userErr := getUserId(c.Param("user_id"))
	if userErr != nil {
		c.JSON(userErr.Status, &userErr)
		return
	}

	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user.ID = userId

	isPartial := c.Request.Method == http.MethodPatch

	result, err2 := services.UserService.UpdateUser(isPartial, user)
	if err2 != nil {
		c.JSON(err2.Status, err2)
		return
	}
	c.JSON(http.StatusOK, result)
}

func Delete(c *gin.Context) {
	// userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	// if userErr != nil {
	// 	err := errors.NewBadRequestError("user id should be a number")
	// 	c.JSON(err.Status, &err)
	// 	return
	// }
	userId, userErr := getUserId(c.Param("user_id"))
	if userErr != nil {
		c.JSON(userErr.Status, &userErr)
		return
	}

	if err := services.UserService.DeleteUser(userId); err != nil {
		c.JSON(err.Status, &err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
	//c.String(http.StatusOK, "deleted")
}

func Search(c *gin.Context) {
	//! localhost:8080/users/search?status=active
	status := c.Query("status")

	users, err := services.UserService.SearchUser(status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	result := users.Marshal(c.GetHeader("X-Public") == "true")
	// result := make([]interface{}, len(users))
	// for index, user := range users {
	// 	result[index] = user.Marshal(c.GetHeader("X-Public") == "true")
	// }

	c.JSON(http.StatusOK, result)

}

func Login(c *gin.Context) {
	var request users.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		rerr := errors.NewBadRequestError("invalid json body")
		c.JSON(rerr.Status, rerr)
		return
	}
	user, err := services.UserService.LoginUser(request)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	result := user.Marshal(c.GetHeader("X-Public") == "true")
	c.JSON(http.StatusOK, result)
}

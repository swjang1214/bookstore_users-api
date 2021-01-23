package app

import (
	"github.com/swjang1214/bookstore_users-api/controllers/ping"
	"github.com/swjang1214/bookstore_users-api/controllers/users_controller"
	_ "github.com/swjang1214/bookstore_users-api/domain/users"
)

func mapUrls() {

	router.GET("/ping", ping.Ping)
	router.POST("/users/login", users_controller.Login)
	// users CRUD
	router.POST("/users", users_controller.Create)
	router.GET("/users/:user_id", users_controller.Get)
	router.PUT("/users/:user_id", users_controller.Update)
	router.PATCH("/users/:user_id", users_controller.Update)
	router.DELETE("/users/:user_id", users_controller.Delete)
	//
	router.GET("/internal/users/search", users_controller.Search)
}

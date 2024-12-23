package handler

import (
	"Batumbuah/modules/v1/utilities/user/repository"
	"Batumbuah/modules/v1/utilities/user/service"
	jwt "Batumbuah/pkg/jwt"

	"gorm.io/gorm"
)

type userHandler struct {
	userService service.Service
	jwtoken     jwt.JwToken
}

func NewUserHandler(userService service.Service, jwtoken jwt.JwToken) *userHandler {
	return &userHandler{userService, jwtoken}
}

func Handler(db *gorm.DB) *userHandler {
	userRepository := repository.NewRepository(db)
	userService := service.NewService(userRepository)
	userHandler := NewUserHandler(userService, jwt.NewJwToken())

	return userHandler
}

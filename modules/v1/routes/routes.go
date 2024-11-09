package routes

import (
	"Batumbuah/app/config"
	"Batumbuah/app/middlewares"

	// deviceHandlerV1 "Batumbuah/modules/v1/utilities/device/handler"
	userHandlerV1 "Batumbuah/modules/v1/utilities/user/handler"
	userViewV1 "Batumbuah/modules/v1/utilities/user/view"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ParseTmpl(router *gin.Engine) *gin.Engine { //Load HTML Template
	router.Static("/assets", "./public/assets")
	router.Static("/img", "./public/assets/img")
	router.Static("/css", "./public/assets/css")
	router.Static("/js", "./public/assets/js")
	router.Static("/vendor", "./public/assets/vendor")
	return router
}

func Init(db *gorm.DB, conf config.Conf, router *gin.Engine) *gin.Engine {
	// deviceHandlerV1 := deviceHandlerV1.Handler(db, conf)
	// deviceViewV1 := deviceviewV1.View(db, conf)
	userHandlerV1 := userHandlerV1.Handler(db)
	userViewV1 := userViewV1.View(db)

	//Routing API Service
	// api := router.Group("/api/v1")

	// Routing Website Service
	user := router.Group("/")
	user.GET("/login", middlewares.LoggedIn(), userViewV1.Login)
	user.POST("/register", userHandlerV1.Register)
	user.POST("/login", userHandlerV1.Login)
	user.GET("/logout", userHandlerV1.Logout)

	user.GET("/", middlewares.IsLogin(), userViewV1.Index)

	router = ParseTmpl(router)
	return router
}

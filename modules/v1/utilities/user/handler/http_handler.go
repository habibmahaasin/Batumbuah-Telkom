package handler

import (
	"Batumbuah/modules/v1/utilities/user/models"
	"Batumbuah/pkg/helpers"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (h *userHandler) Register(c *gin.Context) {
	var input models.RegisterInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Register the user
	err := h.userService.Register(input.FullName, input.Email, input.Password, input.Address, 2)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Registration successful, please log in",
	})
}

func (h *userHandler) Login(c *gin.Context) {
	session := sessions.Default(c)
	var input models.LoginInput

	err := c.ShouldBind(&input)
	if err != nil {
		log.Println(err)
		return
	}

	user, err := h.userService.Login(input)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title":   "Login",
			"message": "Username/ Password yang anda masukkan Salah!",
		})
		return
	}

	token, _ := h.jwtoken.GenerateToken(user.UserID, user.FullName, int(user.RoleID))
	c.SetCookie("Token", token, 21600, "/", "localhost", false, true)

	session.Set("email", user.Email)
	session.Set("full_name", user.FullName)
	session.Set("user_id", user.UserID)
	session.Options(sessions.Options{
		MaxAge: 3600 * 6,
	})
	session.Save()

	c.Redirect(http.StatusFound, "/")
}

func (h *userHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "Batumbuah",
		MaxAge: -1,
	})

	c.Redirect(http.StatusFound, "/login")
}

func (h *userHandler) RegisterPlant(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id").(string)
	email := session.Get("email").(string)

	var input models.RegisterPlantInput
	if err := c.ShouldBind(&input); err != nil {
		helpers.SetFlashMessage(c, "error", err.Error())
        c.Redirect(http.StatusFound, "/")
		return
	}

	err := h.userService.RegisterPlant(userID, input.Name, email)
	if err != nil {
		helpers.SetFlashMessage(c, "error", err.Error())
        c.Redirect(http.StatusFound, "/")
		return
	}

	helpers.SetFlashMessage(c, "success", "Plant registered successfully")
    c.Redirect(http.StatusFound, "/")
}

func (h *userHandler) GetPlantByUserID(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id").(string)

	plants, err := h.userService.GetPlantByUserID(userID)
	if err != nil {
		helpers.SetFlashMessage(c, "error", err.Error())
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   plants,
	})
}

func (h *userHandler) CheckIn(c *gin.Context) {
    session := sessions.Default(c)
    userID := session.Get("user_id").(string)

    plantID := c.Param("id")

    var input models.CheckInInput
    if err := c.ShouldBind(&input); err != nil {
        helpers.SetFlashMessage(c, "error", err.Error())
        c.Redirect(http.StatusFound, "/plant/" + plantID)
        return
    }

    err := h.userService.CheckIn(userID, plantID, "images.jpeg", input.Note)
    if err != nil {
        helpers.SetFlashMessage(c, "error", err.Error())
        c.Redirect(http.StatusFound, "/plant/" + plantID)
        return
    }

    helpers.SetFlashMessage(c, "success", "Check-in successful")
    c.Redirect(http.StatusFound, "/")
}
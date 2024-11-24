package handler

import (
	"Batumbuah/app/config"
	"Batumbuah/modules/v1/utilities/user/models"
	"Batumbuah/pkg/helpers"
	"bytes"
	"fmt"
	"image"
	"log"
	"net/http"
	"time"

	"github.com/disintegration/imaging"
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
	conf, _ := config.Init()

    var input models.CheckInInput
    if err := c.ShouldBind(&input); err != nil {
        helpers.SetFlashMessage(c, "error", err.Error())
        c.Redirect(http.StatusFound, "/plant/"+plantID)
        return
    }

    var buffer []byte
    var contentType string
    var err error
    var imageName string

    imageFile, err := c.FormFile("image")
    if err == nil && imageFile != nil {
        buffer, _, contentType, err = h.readFile(c)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "status":  "error",
                "message": err.Error(),
            })
            return
        }

        isValidContentType := func(contentType string) bool {
            return contentType == "image/jpeg" || contentType == "image/png" || contentType == "image/gif"
        }

        if !isValidContentType(contentType) {
            c.JSON(http.StatusBadRequest, gin.H{
                "status":  "error",
                "message": "Invalid file type.",
            })
            return
        }

        img, _, err := image.Decode(bytes.NewReader(buffer))
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "status":  "error",
                "message": "Failed to decode image.",
            })
            return
        }

        compressedImg := imaging.Resize(img, 400, 0, imaging.Lanczos)
        compressedBuffer := new(bytes.Buffer)
        
        if err := imaging.Encode(compressedBuffer, compressedImg, imaging.JPEG); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "status":  "error",
                "message": "Failed to compress image.",
            })
            return
        }

        buffer = compressedBuffer.Bytes()

        timestamp := time.Now().Format("20060102_150405")
        ext := "jpeg"

        imageName = fmt.Sprintf("%s_%s.%s", plantID, timestamp, ext)
        key := fmt.Sprintf("user-%s/%s/%s", userID, plantID, imageName)

        err = h.uploadImageToS3(conf.Storage.Bucket, key, buffer, "image/jpeg")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "status":  "error",
                "message": err.Error(),
            })
            return
        }
    }

    err = h.userService.CheckIn(userID, plantID, imageName, input.Note)
    if err != nil {
        helpers.SetFlashMessage(c, "error", err.Error())
        c.Redirect(http.StatusFound, "/plant/"+plantID)
        return
    }

    helpers.SetFlashMessage(c, "success", "Check-in successful! Congratulations, you've earned 1 point!")
    c.Redirect(http.StatusFound, "/plant/"+plantID)
}

func (h *userHandler) CreateObject(c *gin.Context) {
    var input models.CreateObjectInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": err.Error(),
        })
        return
    }

    err := h.userService.CreateS3Object(input.BucketName, input.Key, input.Content, "ap-southeast-2")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Object created successfully",
    })
}


func (h *userHandler) readFile(c *gin.Context) ([]byte, string, string, error) {
    file, header, err := c.Request.FormFile("image")
    if err != nil {
        return nil, "", "", fmt.Errorf("failed to upload image: %w", err)
    }
    defer file.Close()

    fileSize := header.Size
    buffer := make([]byte, fileSize)
    if _, err := file.Read(buffer); err != nil {
        return nil, "", "", fmt.Errorf("failed to read file content: %w", err)
    }

    contentType := header.Header.Get("Content-Type")
    return buffer, header.Filename, contentType, nil
}

func (h *userHandler) uploadImageToS3(bucketName, key string, buffer []byte, contentType string) error {
	conf , _ := config.Init()
    err := h.userService.UploadImageToS3(bucketName, key, buffer, contentType, conf.Storage.Region)
    return err
}

func (h *userHandler) UploadImage(c *gin.Context) {
    var input models.UploadImageInput
    if err := c.ShouldBind(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": "Invalid input: " + err.Error(),
        })
        return
    }

    buffer, _, contentType, err := h.readFile(c)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "message": err.Error(),
        })
        return
    }

    err = h.uploadImageToS3(input.BucketName, input.Key, buffer, contentType)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status":  "error",
            "message": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "Image uploaded successfully",
    })
}

func (h *userHandler) UpdatePlantName(c *gin.Context) {
    plantID := c.Param("id")
    var input models.UpdatePlantNameInput

    if err := c.ShouldBind(&input); err != nil {
        helpers.SetFlashMessage(c, "error", err.Error())
        c.Redirect(http.StatusFound, "/")
        return
    }

    err := h.userService.UpdatePlantName(plantID, input.Name)
    if err != nil {
        helpers.SetFlashMessage(c, "error", err.Error())
        c.Redirect(http.StatusFound, "/")
        return
    }

    helpers.SetFlashMessage(c, "success", "Plant name updated successfully")
    c.Redirect(http.StatusFound, "/")
}

func (h *userHandler) DeletePlant(c *gin.Context) {
    plantID := c.Param("id")

    err := h.userService.DeletePlant(plantID)
    if err != nil {
        helpers.SetFlashMessage(c, "error", err.Error())
        c.Redirect(http.StatusFound, "/")
        return
    }

    helpers.SetFlashMessage(c, "success", "Plant deleted successfully")
    c.Redirect(http.StatusFound, "/")
}

func (h *userHandler) UpdatePassword(c *gin.Context) {
    session := sessions.Default(c)
    userID := session.Get("user_id").(string)

    var input models.UpdatePasswordInput
    if err := c.ShouldBind(&input); err != nil {
        helpers.SetFlashMessage(c, "error", err.Error())
        c.Redirect(http.StatusFound, "/profile")
        return
    }

    err := h.userService.UpdatePassword(userID, input.OldPassword, input.NewPassword)
    if err != nil {
        helpers.SetFlashMessage(c, "error", err.Error())
        c.Redirect(http.StatusFound, "/profile")
        return
    }

    helpers.SetFlashMessage(c, "success", "Password updated successfully")
    c.Redirect(http.StatusFound, "/profile")
}
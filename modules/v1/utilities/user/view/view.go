package view

import (
	"Batumbuah/app/config"
	"Batumbuah/pkg/helpers"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (h *userView) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title":   "Login",
		"message": "",
	})
}

func (h *userView) Register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title":   "Register",
	})
}

func (h *userView) Index(c *gin.Context) {
    location, err := time.LoadLocation("Asia/Jakarta")
    if err != nil {
        log.Fatalf("Failed to load time location: %v", err)
    }

    session := sessions.Default(c)
    email := session.Get("email")
    name := session.Get("full_name")
    userID := session.Get("user_id")
    plantLists, _ := h.userService.GetPlantByUserID(userID.(string))
    status, message := helpers.GetAndClearFlashMessage(c)

    for i := range plantLists {
        plantLists[i].DateCreated = plantLists[i].DateCreated.In(location)
        plantLists[i].FormattedDateCreated = plantLists[i].DateCreated.Format("02-01-2006 15:04:05")
    }

    c.HTML(http.StatusOK, "index.html", gin.H{
        "title":       "Home",
        "status":      status,
        "message":     message,
        "email":       email,
        "name":        name,
        "plantLists":  plantLists,
    })
}

func (h *userView) PlantDetail(c *gin.Context) {
    conf , _ := config.Init()
    session := sessions.Default(c)
    plantID := c.Param("id")
    userID := session.Get("user_id")
    status, message := helpers.GetAndClearFlashMessage(c)
    location, err := time.LoadLocation("Asia/Jakarta")
    if err != nil {
        log.Fatalf("Failed to load time location: %v", err)
    }

    plant, _ := h.userService.GetPlantByID(plantID)
    checkInLogs, _ := h.userService.GetCheckInLogs(plantID)
    plantStats, _ := h.userService.GetPlantStatsById(plantID)
    Endpoint := conf.Storage.Endpoint

    for i := range checkInLogs {
        checkInLogs[i].DateCreated = checkInLogs[i].DateCreated.In(location)
        checkInLogs[i].FormattedDateCreated = checkInLogs[i].DateCreated.Format("02-01-2006 15:04:05")
        checkInLogs[i].Image = Endpoint + "/user-" + fmt.Sprintf("%v", userID) + "/" + checkInLogs[i].UserPlantID + "/" + checkInLogs[i].Image
    }

    daysSinceLastCheckIn := 0
    if len(checkInLogs) > 0 {
        lastCheckIn := checkInLogs[len(checkInLogs)-1].DateCreated
        daysSinceLastCheckIn = int(time.Since(lastCheckIn).Hours() / 24)
    }

    c.HTML(http.StatusOK, "plant_detail.html", gin.H{
        "title":               "Plant Detail",
        "plant":               plant,
        "status":              status,
        "message":             message,
        "checkInLogs":         checkInLogs,
        "sumCheckInLogs":      len(checkInLogs),
        "daysSinceLastCheckIn": daysSinceLastCheckIn,
        "plantStats":          plantStats,
        "userID":              userID,
        "plantID":             plantID,
    })
}

func (h *userView) Profile(c *gin.Context) {
    session := sessions.Default(c)
    email := session.Get("email")
    name := session.Get("full_name")
    status, message := helpers.GetAndClearFlashMessage(c)

    c.HTML(http.StatusOK, "profile.html", gin.H{
        "title":       "Profile",
        "status":      status,
        "message":     message,
        "email":       email,
        "name":        name,
    })
}

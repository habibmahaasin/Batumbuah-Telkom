package view

import (
	"Batumbuah/pkg/helpers"
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
    session := sessions.Default(c)
    email := session.Get("email")
    name := session.Get("full_name")
    userID := session.Get("user_id")
    plantLists, _ := h.userService.GetPlantByUserID(userID.(string))
    status, message := helpers.GetAndClearFlashMessage(c)

    for i := range plantLists {
    plantLists[i].DateCreated = plantLists[i].DateCreated.In(time.Local)
    plantLists[i].FormattedDateCreated = plantLists[i].DateCreated.Format("2006-01-02 15:04:05")
}

    c.HTML(http.StatusOK, "index.html", gin.H{
        "title":       "Index",
        "status":      status,
        "message":     message,
        "email":       email,
        "name":        name,
        "plantLists":  plantLists,
    })
}

func (h *userView) PlantDetail(c *gin.Context) {
    plantID := c.Param("id")
    status, message := helpers.GetAndClearFlashMessage(c)

    plant, _ := h.userService.GetPlantByID(plantID)
    checkInLogs, _ := h.userService.GetCheckInLogs(plantID)
    plantStats, _ := h.userService.GetPlantStatsById(plantID)
    
  for i := range checkInLogs {
    checkInLogs[i].DateCreated = checkInLogs[i].DateCreated.In(time.Local)
    checkInLogs[i].FormattedDateCreated = checkInLogs[i].DateCreated.Format("2006-01-02 15:04:05")
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

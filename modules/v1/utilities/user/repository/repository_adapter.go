package repository

import (
	"Batumbuah/modules/v1/utilities/user/models"

	"gorm.io/gorm"
)

type Repository interface {
	GetUserByEmail(email string) (models.User, error)
	CreateUser(user *models.User) error
	RegisterPlant(plant *models.UserPlant, plantStats *models.PlantStats, testInfo *models.TestInformation) error
	GetPlantByUserID(userID string) ([]models.UserPlant, error)
	GetPlantByID(plantID string) (models.UserPlant, error)
	PlantCheckIn(UserPlantID, image, note string) error
	GetLastCheckInTime(UserPlantID string) (models.CheckInLog, error)
	GetCheckInLogs(UserPlantID string) ([]models.CheckInLog, error)
	GetPlantStatsById(plantID string) (models.PlantStats, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

package service

import (
	"Batumbuah/modules/v1/utilities/user/models"
	"Batumbuah/modules/v1/utilities/user/repository"
)

type Service interface {
	Login(input models.LoginInput) (models.User, error)
	Register(fullName, email, password, address string, roleID int) error
	RegisterPlant(userID, name, email string) error
	GetPlantByUserID(userID string) ([]models.UserPlant, error)
	GetPlantByID(plantID string) (models.UserPlant, error)
	CheckIn(userID, plantID, image, note string) error 
	GetCheckInLogs(UserPlantID string) ([]models.CheckInLog, error)
	GetPlantStatsById(plantID string) (models.PlantStats, error)
	CreateS3Object(bucketName, key, content, region string) error
	UploadImageToS3(bucketName, key string, content []byte, contentType, region string) error
	UpdatePlantName(PlantID string, name string) error
	DeletePlant(PlantID string) error
	UpdatePassword(userID, oldPassword string, newPassword string) error
}

type service struct {
	repository repository.Repository
}

func NewService(repository repository.Repository) *service {
	return &service{repository}
}

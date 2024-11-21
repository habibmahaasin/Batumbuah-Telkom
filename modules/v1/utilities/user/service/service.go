package service

import (
	"Batumbuah/modules/v1/utilities/user/models"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (n *service) Login(input models.LoginInput) (models.User, error) {
	email := input.Email
	password := input.Password
	user, _ := n.repository.GetUserByEmail(email)

	if user.UserID == "" {
		return user, errors.New("email yang Anda masukan salah")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println(err)
		return user, errors.New("password yang Anda masukan salah/tidak terdaftar")
	}
	return user, nil
}

func (s *service) Register(fullName, email, password, address string, roleID int) error {

	existingUser, err := s.repository.GetUserByEmail(email)
	if err != nil {
		return err
	}

	if existingUser.UserID != "" {
		return errors.New("email already exists")
	}

	userID := uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		UserID:   userID,
		FullName: fullName,
		Email:    email,
		Password: string(hashedPassword),
		Address:  address,
		RoleID:   int64(roleID),
		DateCreated:        time.Now(),
        DateUpdated:        time.Now(),
	}

	return s.repository.CreateUser(user)
}

func (s *service) RegisterPlant(userID, name, email string) error {
    userPlant := &models.UserPlant{
        UserID:  userID,
        PlantID: uuid.New().String(),
        Name:    name,
        DateCreated: time.Now(),
        DateUpdated: time.Now(),
    }

    plantStats := &models.PlantStats{
        PlantID:            userPlant.PlantID,
        TotalCheckIn:       0,
        IsPreTested:        false,
        IsPostTested:       false,
        IsAvailableToRedeem: false,
        IsRedeemedReward:   false,
        DateCreated:        time.Now(),
        DateUpdated:        time.Now(),
    }

    testInfo := &models.TestInformation{
        PlantID:    userPlant.PlantID,
        Email:      email,
        DateCreated: time.Now(),
        DateUpdated: time.Now(),
    }

    return s.repository.RegisterPlant(userPlant, plantStats, testInfo)
}

func (s *service) GetPlantByUserID(userID string) ([]models.UserPlant, error) {
	return s.repository.GetPlantByUserID(userID)
}

func (s *service) GetPlantByID(plantID string) (models.UserPlant, error) {
	return s.repository.GetPlantByID(plantID)
}

func (s *service) PlantCheckIn(UserPlantID, image, note string) error {
	return s.repository.PlantCheckIn(UserPlantID, image, note)
}

func (s *service) CheckIn(userID, plantID, image, note string) error {
    if plantID == "" {
        return errors.New("plant ID cannot be empty")
    }

    lastCheckIn, err := s.repository.GetLastCheckInTime(plantID)
    if err != nil {
        if err.Error() == "no check-in found for the plant" {
            return s.repository.PlantCheckIn(plantID, image, note)
        }
        return err 
    }

	checkInRule, err := s.repository.GetCheckInRule()
	if err != nil {
        return err
    }
	
    allowedInterval := time.Duration(checkInRule.Range) * 24 * time.Hour

    if time.Since(lastCheckIn.DateCreated) < allowedInterval {
        return fmt.Errorf("check-in allowed only once every %d days", checkInRule.Range)
    }

    return s.repository.PlantCheckIn(plantID, image, note)
}


func (s *service) GetCheckInLogs(UserPlantID string) ([]models.CheckInLog, error) {
	return s.repository.GetCheckInLogs(UserPlantID)
}

func (s *service) GetPlantStatsById(plantID string) (models.PlantStats, error) {
	return s.repository.GetPlantStatsById(plantID)
}

// func (s *service) GetUserStats(userID string) (models.UserStats, error) {
// 	return s.repository.GetUserStats(userID)
// }

// func (s *service) UpdatePreTestStatus(userID string, email string, status bool) error {
// 	err := s.repository.UpdateTestInformation(models.TestInformation{
// 		UserID: userID,
// 		Email:  email,
// 	})
	
// 	if err != nil {
// 		return fmt.Errorf("failed to create test information: %w", err)
// 	}

// 	err = s.repository.UpdateUserStats(userID, models.UserStats{
// 		UserID:      userID,
// 		IsPreTested: status,
// 		DateUpdated: time.Now(),
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to update user stats: %w", err)
// 	}

// 	return nil
// }

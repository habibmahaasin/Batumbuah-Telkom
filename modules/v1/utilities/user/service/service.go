package service

import (
	appConfig "Batumbuah/app/config"
	"Batumbuah/modules/v1/utilities/user/models"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
		return user, errors.New("password yang Anda masukan salah/tidak terdaftar")
	}
	return user, nil
}

func (s *service) Register(fullName, email, password, address string, roleID int) error {
	conf , _ := appConfig.Init()
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
		UserID:      userID,
		FullName:    fullName,
		Email:       email,
		Password:    string(hashedPassword),
		Address:     address,
		RoleID:      int64(roleID),
		DateCreated: time.Now(),
		DateUpdated: time.Now(),
	}

	if err := s.repository.CreateUser(user); err != nil {
		return err
	}

	bucketName := conf.Storage.Bucket
	key := fmt.Sprintf("user-%s/", userID)
	content := "Initial content for user folder"
	region := conf.Storage.Region

	if err := s.CreateS3Object(bucketName, key, content, region); err != nil {
		return fmt.Errorf("failed to create S3 bucket: %w", err)
	}

	return nil
}


func (s *service) RegisterPlant(userID, name, email string) error {
	conf , _ := appConfig.Init()
	plantID := uuid.New().String()
    userPlant := &models.UserPlant{
        UserID:  userID,
        PlantID: plantID,
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

	bucketName := conf.Storage.Bucket
	key := fmt.Sprintf("user-%s/%s/", userID, plantID)
	content := "Initial content for user plant"
	region := conf.Storage.Region

	if err := s.CreateS3Object(bucketName, key, content, region); err != nil {
		return fmt.Errorf("failed to create S3 bucket: %w", err)
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
			if updateErr := s.repository.UpdatePreTestStatus(plantID, true); updateErr != nil {
				return updateErr
			}
			return s.repository.PlantCheckIn(plantID, image, note)
		}
		return err
	}

	plantStats , err := s.repository.GetPlantStatsById(plantID)
	if err != nil {
		return err
	}

	if plantStats.TotalCheckIn == 4 {
		s.repository.UpdateRedeemRewardStatus(plantID, true)
		s.repository.UpdatePostTestStatus(plantID, true)
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

func (s *service) CreateS3Object(bucketName, key, content, region string) error {
	conf , _ := appConfig.Init()
    staticCreds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
       conf.Storage.Access_Key, conf.Storage.Private_Key,
        "",
    ))

    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithRegion(region),
        config.WithCredentialsProvider(staticCreds),
    )
    if err != nil {
        return err
    }

    client := s3.NewFromConfig(cfg)

    var contentReader io.Reader
    if content == "" {
        key = strings.TrimSuffix(key, "/") + "/"
        contentReader = bytes.NewReader([]byte{})
    } else {
        contentReader = strings.NewReader(content)
    }

    _, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(key),
        Body:   contentReader,
    })
    if err != nil {
        return err
    }

    log.Printf("Object %s created successfully in bucket %s", key, bucketName)
    return nil
}

func (s *service) UploadImageToS3(bucketName, key string, content []byte, contentType, region string) error {
	conf , _ := appConfig.Init()

    staticCreds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
       conf.Storage.Access_Key, conf.Storage.Private_Key,
        "",
    ))

    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithRegion(region),
        config.WithCredentialsProvider(staticCreds),
    )
    if err != nil {
        return err
    }

    client := s3.NewFromConfig(cfg)

    _, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
        Bucket:      aws.String(bucketName),
        Key:         aws.String(key),
        Body:        bytes.NewReader(content),
        ContentType: aws.String(contentType),
    })
    if err != nil {
        return err
    }

    log.Printf("Image %s uploaded successfully to bucket %s", key, bucketName)
    return nil
}

func (s *service) UpdatePlantName(PlantID string, name string) error {
	return s.repository.UpdatePlantName(PlantID, name)
}

func (s *service) DeletePlant(PlantID string) error {
	return s.repository.DeletePlant(PlantID)
}

func (s *service) UpdatePassword(userID, oldPassword, newPassword string) error {
   	isValidPassword := func(password string) bool {
    return len(password) >= 8 &&
        regexp.MustCompile(`[A-Z]`).MatchString(password) &&
        regexp.MustCompile(`\d`).MatchString(password)
	}

    if !isValidPassword(newPassword) {
        return fmt.Errorf("password format is invalid, please refer to the password policy")
    }


    currentPasswordHash, err := s.repository.GetPasswordHash(userID)
    if err != nil {
        return fmt.Errorf("failed to retrieve current password: %v", err)
    }

    err = bcrypt.CompareHashAndPassword([]byte(currentPasswordHash), []byte(oldPassword))
    if err != nil {
        return fmt.Errorf("old password does not match")
    }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("failed to hash new password: %v", err)
    }

    err = s.repository.UpdatePassword(userID, string(hashedPassword))
    if err != nil {
        return fmt.Errorf("failed to update password: %v", err)
    }

    return nil
}

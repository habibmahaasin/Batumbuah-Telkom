package repository

import (
	"Batumbuah/modules/v1/utilities/user/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

func (r *repository) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	err := r.db.Where("email = ?", email).Find(&user).Error
	return user, err
}

func (r *repository) CreateUser(user *models.User) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(user).Error; err != nil {
            return err
        }
        return nil
    })
}

func (r *repository) RegisterPlant(plant *models.UserPlant, plantStats *models.PlantStats, testInfo *models.TestInformation) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(plant).Error; err != nil {
            return err
        }

        plantStats.PlantID = plant.PlantID
        testInfo.PlantID = plant.PlantID

        if err := tx.Create(plantStats).Error; err != nil {
            return err
        }

        if err := tx.Create(testInfo).Error; err != nil {
            return err
        }

        return nil
    })
}
func (r *repository) GetPlantByUserID(userID string) ([]models.UserPlant, error) {
    var userPlant []models.UserPlant
    err := r.db.
        Where("user_id = ?", userID).
        Preload("PlantStats").
        Order("date_created DESC").
        Find(&userPlant).Error
    return userPlant, err
}

func (r *repository) GetPlantByID(plantID string) (models.UserPlant, error) {
    var userPlant models.UserPlant
    err := r.db.Where("plant_id = ?", plantID).First(&userPlant).Error
    return userPlant, err
}

func (r *repository) PlantCheckIn(UserPlantID, image, note string) error {
    var userPlantExists bool
    if err := r.db.Model(&models.UserPlant{}).
        Where("plant_id = ?", UserPlantID).
        Select("count(*) > 0").
        Scan(&userPlantExists).Error; err != nil {
        return err
    }

    if !userPlantExists {
        return errors.New("UserPlantID not found")
    }

    return r.db.Transaction(func(tx *gorm.DB) error {
        checkInLog := models.CheckInLog{
            UserPlantID: UserPlantID,
            Image:       image,
            Note:        note,
            DateCreated: time.Now(),
            DateUpdated: time.Now(),
        }

        if err := tx.Create(&checkInLog).Error; err != nil {
            return err
        }

        var checkInCount int64
        if err := tx.Model(&models.CheckInLog{}).
            Where("plant_id = ?", UserPlantID).
            Count(&checkInCount).Error; err != nil {
            return err
        }

        if err := tx.Model(&models.PlantStats{}).
            Where("plant_id = ?", UserPlantID).
            Update("total_check_in", checkInCount).Error; err != nil {
            return err
        }

        return nil
    })
}

func (r *repository) GetLastCheckInTime(plantID string) (models.CheckInLog, error) {
    var checkInLog models.CheckInLog
    if plantID == "" {
        return checkInLog, errors.New("plant ID cannot be empty")
    }

    err := r.db.Where("plant_id = ?", plantID).Order("date_created desc").First(&checkInLog).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return checkInLog, errors.New("no check-in found for the plant")
        }
        return checkInLog, err
    }
    return checkInLog, nil
}

func (r *repository) GetCheckInLogs(UserPlantID string) ([]models.CheckInLog, error) {
    var checkInLogs []models.CheckInLog
    err := r.db.Where("plant_id = ?", UserPlantID).Find(&checkInLogs).Error
    return checkInLogs, err
}

func (r *repository) GetPlantStatsById(plantID string) (models.PlantStats, error) {
    var plantStats models.PlantStats
    err := r.db.Where("plant_id = ?", plantID).First(&plantStats).Error
    return plantStats, err
}

func (r *repository) GetCheckInRule() (models.CheckInRule, error) {
    var checkInRule models.CheckInRule
    err := r.db.First(&checkInRule).Error
    return checkInRule, err
}

func (r *repository) UpdatePreTestStatus(plantID string, status bool) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Model(&models.PlantStats{}).
            Where("plant_id = ?", plantID).
            Update("is_pre_tested", status).Error; err != nil {
            return err
        }

        return nil
    })
}

func (r *repository) UpdatePostTestStatus(PlantID string , status bool) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Model(&models.PlantStats{}).
            Where("plant_id = ?", PlantID).
            Update("is_post_tested", status).Error; err != nil {
            return err
        }

        return nil
    })
} 

func (r *repository) UpdateRedeemRewardStatus(PlantID string , status bool) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Model(&models.PlantStats{}).
            Where("plant_id = ?", PlantID).
            Update("is_available_to_redeem", status).Error; err != nil {
            return err
        }

        return nil
    })
}

func (r *repository) UpdatePlantName(PlantID string, name string) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Model(&models.UserPlant{}).
            Where("plant_id = ?", PlantID).
            Update("name", name).Error; err != nil {
            return err
        }

        return nil
    })
}

func (r *repository) DeletePlant(PlantID string) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Where("plant_id = ?", PlantID).Delete(&models.UserPlant{}).Error; err != nil {
            return err
        }

        return nil
    })
}

func (r *repository) GetPasswordHash(userID string) (string, error) {
    var user models.User
    err := r.db.Where("user_id = ?", userID).First(&user).Error
    return user.Password, err
}

func (r *repository) UpdatePassword(userID string, newPassword string) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Model(&models.User{}).
            Where("user_id = ?", userID).
            Update("password", newPassword).Error; err != nil {
            return err
            }

        return nil
    })
}
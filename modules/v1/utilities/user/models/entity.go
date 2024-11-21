package models

import "time"

// User model definition
type User struct {
    UserID      string       `gorm:"primaryKey;column:user_id"`
    FullName    string       `gorm:"column:full_name"`
    Email       string       `gorm:"column:email"`
    RoleID      int64        `gorm:"column:role_id"`
    Password    string       `gorm:"column:password"`
    Address     string       `gorm:"column:address"`
    DateCreated time.Time    `gorm:"column:date_created"`
    DateUpdated time.Time    `gorm:"column:date_updated"`
    UserPlants  []UserPlant  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

// UserPlant model definition
type UserPlant struct {
    PlantID         string          `gorm:"primaryKey;column:plant_id"`
    UserID          string          `gorm:"column:user_id"` // References User.UserID
    Name            string          `gorm:"column:name"`
    DateCreated time.Time `gorm:"column:date_created" json:"date_created"`
    DateUpdated time.Time `gorm:"column:date_updated" json:"date_updated"`
    CheckInLogs     []CheckInLog    `gorm:"foreignKey:UserPlantID;references:PlantID;constraint:OnDelete:CASCADE;"`
    PlantStats      PlantStats      `gorm:"foreignKey:PlantID;references:PlantID;constraint:OnDelete:CASCADE;"`
    TestInformation TestInformation `gorm:"foreignKey:PlantID;references:PlantID;constraint:OnDelete:CASCADE;"`
}

// CheckInLog model definition
type CheckInLog struct {
    ID          uint      `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
    UserPlantID string    `gorm:"column:plant_id" json:"plant_id"` // References UserPlant.PlantID
    Image       string    `gorm:"column:image" json:"image"`
    Note        string    `gorm:"column:note" json:"note"`
    DateCreated time.Time `gorm:"column:date_created" json:"date_created"`
    DateUpdated time.Time `gorm:"column:date_updated" json:"date_updated"`
    FormattedDateCreated string `gorm:"-" json:"formatted_date_created"`
}

// PlantStats model definition
type PlantStats struct {
    PlantID             string    `gorm:"primaryKey;column:plant_id"` // References UserPlant.PlantID
    TotalCheckIn        int       `gorm:"column:total_check_in" json:"total_check_in"`
    IsPreTested         bool      `gorm:"column:is_pre_tested" json:"is_pre_tested"`
    IsPostTested        bool      `gorm:"column:is_post_tested" json:"is_post_tested"`
    IsAvailableToRedeem bool      `gorm:"column:is_available_to_redeem" json:"is_available_to_redeem"`
    IsRedeemedReward    bool      `gorm:"column:is_redeemed_reward" json:"is_redeemed_reward"`
    DateCreated         time.Time `gorm:"column:date_created" json:"date_created"`
    DateUpdated         time.Time `gorm:"column:date_updated" json:"date_updated"`
}

// TestInformation model definition
type TestInformation struct {
    PlantID     string    `gorm:"primaryKey;column:plant_id" json:"plant_id"` // References UserPlant.PlantID
    Email       string    `gorm:"column:email" json:"email"`
    DateCreated time.Time `gorm:"column:date_created" json:"date_created"`
    DateUpdated time.Time `gorm:"column:date_updated" json:"date_updated"`
}

type CheckInRule struct {
    ID          uint      `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
    Name       string    `gorm:"column:name" json:"name"`
    Range       int       `gorm:"column:range" json:"range"`
    DateCreated time.Time `gorm:"column:date_created" json:"date_created"`
    DateUpdated time.Time `gorm:"column:date_updated" json:"date_updated"`
}
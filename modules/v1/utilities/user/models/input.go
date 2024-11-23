package models

type LoginInput struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type RegisterInput struct {
	FullName string `json:"full_name" form:"full_name" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
	Address  string `json:"address" form:"address"`
}

type CheckInInput struct {
	// Image string `json:"image" form:"image"`
	Note  string `json:"note" form:"note"`
}

type PreTestStatusInput struct {
	Email  string `json:"email" form:"email" binding:"required,email"`
	Status bool   `json:"status" form:"status" binding:"required"`
}

type RegisterPlantInput struct {
	Name string `json:"name" form:"name" binding:"required"`
}

type CreateObjectInput struct {
    BucketName string `json:"bucketName" binding:"required"` // Name of the bucket
    Key        string `json:"key" binding:"required"`        // Key for the file/folder (path)
    Content    string `json:"content"`                       // Content of the file (empty for folder)
}

type UploadImageInput struct {
    BucketName string `form:"bucketName" binding:"required"` // Name of the target bucket
    Key        string `form:"key" binding:"required"`        // S3 object key (path/filename)
}
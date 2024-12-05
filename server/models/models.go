package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Username     string               `bson:"username" json:"username"`
	Email        string               `bson:"email" json:"email"`
	PasswordHash string               `bson:"password_hash" json:"password_hash"`
	Permissions  string               `bson:"permissions" json:"permissions"`
	CreatedAt    time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time            `bson:"updated_at" json:"updated_at"`
	LastLoginAt  time.Time            `bson:"last_login_at" json:"last_login_at"`
	Payments     []Payment            `bson:"payments" json:"payments"`
	Jobs         []primitive.ObjectID `bson:"jobs" json:"jobs"`
}

type Payment struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Price         string             `bson:"price" json:"price"`
	PaymentMethod string             `bson:"payment_method" json:"payment_method"`
	PaymentStatus string             `bson:"payment_status" json:"payment_status"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
	JobID         primitive.ObjectID `bson:"job_id" json:"job_id"`
}

type Job struct {
	ID                      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID                  primitive.ObjectID `bson:"user_id" json:"user_id"`
	Title                   string             `bson:"title" json:"title"`
	Status                  string             `bson:"status" json:"status"`
	SourceLanguage          string             `bson:"source_language" json:"source_language"`
	FileFormat              string             `bson:"file_format" json:"file_format"`
	Description             string             `bson:"description" json:"description"`
	InputFile               string             `bson:"input_file" json:"input_file"`
	OutputFile              string             `bson:"output_file" json:"output_file"`
	CreatedAt               time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt               time.Time          `bson:"updated_at" json:"updated_at"`
	EstimatedFinishDatetime time.Time          `bson:"estimated_finish_datetime" json:"estimated_finish_datetime"`
	HostID                  primitive.ObjectID `bson:"host_id" json:"host_id"`
}

type Server struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Hostname      string               `bson:"hostname" json:"hostname"`
	Address       string               `bson:"address" json:"address"`
	Description   string               `bson:"description" json:"description"`
	Status        string               `bson:"status" json:"status"`
	CreatedAt     time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time            `bson:"updated_at" json:"updated_at"`
	CurrentJobs   []primitive.ObjectID `bson:"current_jobs" json:"current_jobs"`
	CompletedJobs []primitive.ObjectID `bson:"completed_jobs" json:"completed_jobs"`
	CPUInfo       string               `bson:"cpu_info" json:"cpu_info"`
	GPUInfo       string               `bson:"gpu_info" json:"gpu_info"`
	RAMSizeGB     int32                `bson:"ram_size_gb" json:"ram_size_gb"`
}

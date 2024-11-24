package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty"`
	Username     string               `bson:"username"`
	Email        string               `bson:"email"`
	PasswordHash string               `bson:"password_hash"`
	Permissions  string               `bson:"permissions"`
	CreatedAt    time.Time            `bson:"created_at"`
	UpdatedAt    time.Time            `bson:"updated_at"`
	LastLoginAt  time.Time            `bson:"last_login_at"`
	Payments     []Payment            `bson:"payments"`
	Jobs         []primitive.ObjectID `bson:"jobs"`
}

type Payment struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Price         string             `bson:"price"`
	PaymentMethod string             `bson:"payment_method"`
	PaymentStatus string             `bson:"payment_status"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at"`
}

type Job struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	UserID         primitive.ObjectID `bson:"user_id"`
	Title          string             `bson:"title"`
	Status         string             `bson:"status"`
	SourceLanguage string             `bson:"source_language"`
	FileFormat     string             `bson:"file_format"`
	Description    string             `bson:"description"`
	InputFile      string             `bson:"input_file"`
	OutputFile     string             `bson:"output_file"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
	HostID         primitive.ObjectID `bson:"host_id"`
}

type Server struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty"`
	Hostname       string               `bson:"hostname"`
	Address        string               `bson:"address"`
	Description    string               `bson:"description"`
	Status         string               `bson:"status"`
	CreatedAt      time.Time            `bson:"created_at"`
	UpdatedAt      time.Time            `bson:"updated_at"`
	CurrentTasks   []primitive.ObjectID `bson:"current_tasks"`
	CompletedTasks []primitive.ObjectID `bson:"completed_tasks"`
}

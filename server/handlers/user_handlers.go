package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	"github.com/moevm/nosql2h24-transcribtion/db"
	"github.com/moevm/nosql2h24-transcribtion/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All good!"))
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser models.User

	if err := render.DecodeJSON(r.Body, &newUser); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	fmt.Println(newUser)

	if newUser.Username == "" || newUser.Email == "" || newUser.PasswordHash == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.PasswordHash), bcrypt.DefaultCost)
	//if err != nil {
	//	http.Error(w, "Error hashing password", http.StatusInternalServerError)
	//	return
	//}
	//newUser.PasswordHash = string(hashedPassword)

	newUser.ID = primitive.NewObjectID()
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()
	newUser.Permissions = "user"

	usersCollection := db.GetCollection("users")

	_, err := usersCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		http.Error(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	newUser.PasswordHash = ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

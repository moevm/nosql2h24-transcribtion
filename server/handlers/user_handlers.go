package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/moevm/nosql2h24-transcribtion/db"
	"github.com/moevm/nosql2h24-transcribtion/models"
	"github.com/moevm/nosql2h24-transcribtion/scheduler"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
	"time"
)

// GetUsers обрабатывает GET-запросы для получения списка пользователей с возможностью фильтрации и пагинации.
// Параметры фильтрации:
// - username (опционально): фильтрация по имени пользователя (например, ?username=John)
// - email (опционально): фильтрация по email (например, ?email=example@example.com)
// - created_after (опционально): фильтрация по дате создания, только пользователи, созданные после указанной даты (формат: YYYY-MM-DD)
// - status (опционально): фильтрация по статусу аккаунта (например, ?status=active)
// Параметры пагинации:
// - page (опционально): номер страницы для пагинации (по умолчанию 1)
// - page_size (опционально): количество пользователей на странице (по умолчанию 10)
// Ответ:
// Возвращает список пользователей в формате JSON, соответствующих фильтрам и с учетом пагинации.
// Если параметры запроса некорректны, возвращает ошибку с кодом 400 (Bad Request).
func GetUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	filter := bson.M{}

	// Фильтрация по имени пользователя (например, ?username=John)
	if username := queryParams.Get("username"); username != "" {
		filter["username"] = username
	}

	// Фильтрация по email (например, ?email=example@example.com)
	if email := queryParams.Get("email"); email != "" {
		filter["email"] = email
	}

	// Фильтрация по дате создания (например, ?created_after=2024-01-01)
	if createdAfter := queryParams.Get("created_after"); createdAfter != "" {
		createdAfterTime, err := time.Parse("2006-01-02", createdAfter)
		if err != nil {
			http.Error(w, "Invalid created_after date format", http.StatusBadRequest)
			return
		}
		filter["created_at"] = bson.M{"$gte": createdAfterTime}
	}

	// Фильтрация по статусу аккаунта (например, ?status=active)
	if status := queryParams.Get("status"); status != "" {
		filter["status"] = status
	}

	usersCollection := db.GetCollection("users")

	page := queryParams.Get("page")
	pageSize := queryParams.Get("page_size")
	var pageNum, pageSizeInt int64

	if page == "" {
		pageNum = 1
	} else {
		var err error
		pageNum, err = strconv.ParseInt(page, 10, 64)
		if err != nil || pageNum <= 0 {
			http.Error(w, "Invalid page parameter", http.StatusBadRequest)
			return
		}
	}
	if pageSize == "" {
		pageSizeInt = 10
	} else {
		var err error
		pageSizeInt, err = strconv.ParseInt(pageSize, 10, 64)
		if err != nil || pageSizeInt <= 0 {
			http.Error(w, "Invalid page_size parameter", http.StatusBadRequest)
			return
		}
	}

	// Опции пагинации
	opt := options.Find().SetSkip((pageNum - 1) * pageSizeInt).SetLimit(pageSizeInt)

	cursor, err := usersCollection.Find(context.Background(), filter, opt)
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var users []models.User
	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			http.Error(w, "Error decoding user", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, "Error iterating over users", http.StatusInternalServerError)
		return
	}
	for i := range users {
		users[i].PasswordHash = ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// GetUserByID получает пользователя по его уникальному ID.
// Параметры:
//   - id (URL параметр): ID пользователя, который необходимо получить.
//
// Ответ:
//   - 200 OK: Возвращает информацию о пользователе (без пароля).
//   - 400 Bad Request: Неверный формат ID пользователя.
//   - 404 Not Found: Пользователь с таким ID не найден.
//   - 500 Internal Server Error: Ошибка при запросе к базе данных.
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}
	filter := bson.M{"_id": objectID}

	var user models.User
	err = db.GetCollection("users").FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching user", http.StatusInternalServerError)
		}
		return
	}
	user.PasswordHash = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser обновляет информацию о пользователе по его уникальному ID.
// Параметры:
//   - id (URL параметр): ID пользователя, которого нужно обновить.
//
// Тело запроса:
//   - username (string): Новое имя пользователя (опционально).
//   - email (string): Новый email пользователя (опционально).
//   - permissions (string): Новая роль/права пользователя (опционально).
//
// Ответ:
//   - 200 OK: Возвращает обновленную информацию о пользователе (без пароля).
//   - 400 Bad Request: Ошибка при обновлении (например, отсутствуют поля для обновления).
//   - 404 Not Found: Пользователь с таким ID не найден.
//   - 500 Internal Server Error: Ошибка при обновлении пользователя или запросе к базе данных.
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	var updatedUser models.User
	if err := render.DecodeJSON(r.Body, &updatedUser); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if updatedUser.Username == "" && updatedUser.Email == "" && updatedUser.Permissions == "" {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}
	filter := bson.M{"_id": objectID}

	update := bson.M{
		"$set": bson.M{
			"username":    updatedUser.Username,
			"email":       updatedUser.Email,
			"permissions": updatedUser.Permissions,
			"updated_at":  time.Now(),
		},
	}

	_, err = db.GetCollection("users").UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	var user models.User
	err = db.GetCollection("users").FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		http.Error(w, "Error fetching updated user", http.StatusInternalServerError)
		return
	}

	user.PasswordHash = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// PatchUser частично обновляет информацию о пользователе по его уникальному ID.
// Параметры:
//   - id (URL параметр): ID пользователя, которого нужно обновить.
//
// Тело запроса:
//   - username (string): Новое имя пользователя (опционально).
//   - email (string): Новый email пользователя (опционально).
//   - permissions (string): Новая роль/права пользователя (опционально).
//
// Ответ:
//   - 200 OK: Возвращает обновленную информацию о пользователе (без пароля).
//   - 400 Bad Request: Ошибка при обновлении (например, отсутствуют поля для обновления).
//   - 404 Not Found: Пользователь с таким ID не найден.
//   - 500 Internal Server Error: Ошибка при обновлении пользователя или запросе к базе данных.
func PatchUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}
	var patchUser models.User
	if err := render.DecodeJSON(r.Body, &patchUser); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	filter := bson.M{"_id": objectID}

	updateFields := bson.M{}

	if patchUser.Username != "" {
		updateFields["username"] = patchUser.Username
	}
	if patchUser.Email != "" {
		updateFields["email"] = patchUser.Email
	}
	if patchUser.Permissions != "" {
		updateFields["permissions"] = patchUser.Permissions
	}

	if len(updateFields) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}
	updateFields["updated_at"] = time.Now()

	update := bson.M{"$set": updateFields}
	_, err = db.GetCollection("users").UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	var user models.User
	err = db.GetCollection("users").FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		http.Error(w, "Error fetching updated user", http.StatusInternalServerError)
		return
	}
	user.PasswordHash = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
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

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	usersCollection := db.GetCollection("users")
	_, err = usersCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetUserJobs(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	usersCollection := db.GetCollection("users")
	var user models.User
	err = usersCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	jobsCollection := db.GetCollection("jobs")
	var jobs []models.Job
	cursor, err := jobsCollection.Find(context.Background(), bson.M{"_id": bson.M{"$in": user.Jobs}})
	if err != nil {
		http.Error(w, "Error fetching jobs", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())
	if err = cursor.All(context.Background(), &jobs); err != nil {
		http.Error(w, "Error parsing jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func AddUserJob(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var job models.Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	job.ID = primitive.NewObjectID()
	job.UserID = id
	job.CreatedAt = time.Now()
	job.UpdatedAt = job.CreatedAt

	serversCollection := db.GetCollection("servers")
	jobsCollection := db.GetCollection("jobs")
	usersCollection := db.GetCollection("users")

	servers, err := schedul.GetServers(serversCollection)
	if err != nil {
		http.Error(w, "Error fetching servers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	selectedServer, err := schedul.SelectServerWithMinJobs(servers)
	if err != nil {
		http.Error(w, "Error selecting server: "+err.Error(), http.StatusInternalServerError)
		return
	}

	job.HostID = selectedServer.ID
	if err := schedul.AddJobToServer(serversCollection, selectedServer.ID, job.ID); err != nil {
		http.Error(w, "Error adding job to server: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = jobsCollection.InsertOne(context.Background(), job)
	if err != nil {
		http.Error(w, "Error saving job", http.StatusInternalServerError)
		return
	}

	_, err = usersCollection.UpdateOne(context.Background(),
		bson.M{"_id": id},
		bson.M{"$push": bson.M{"jobs": job.ID}},
	)
	if err != nil {
		http.Error(w, "Error updating user jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}

func DeleteUserJob(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	jobID := chi.URLParam(r, "jobId")

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	jobObjectID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	jobsCollection := db.GetCollection("jobs")
	_, err = jobsCollection.DeleteOne(context.Background(), bson.M{"_id": jobObjectID})
	if err != nil {
		http.Error(w, "Error deleting job", http.StatusInternalServerError)
		return
	}

	usersCollection := db.GetCollection("users")
	_, err = usersCollection.UpdateOne(context.Background(),
		bson.M{"_id": userObjectID},
		bson.M{"$pull": bson.M{"jobs": jobObjectID}},
	)
	if err != nil {
		http.Error(w, "Error updating user jobs", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

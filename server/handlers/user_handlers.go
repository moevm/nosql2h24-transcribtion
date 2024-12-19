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
	"log"
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

//1. Получение всех пользователей без фильтров:
//GET /users
//2. Фильтрация по имени пользователя:
//GET /users?username=John
//3. Фильтрация по email:
//GET /users?email=example@example.com
//4. Фильтрация по дате создания (пользователи, созданные после указанной даты):
//GET /users?created_after=2024-01-01
//5. Фильтрация по статусу аккаунта:
//GET /users?status=active
//6. Комбинация фильтров (например, по имени пользователя и email):
//GET /users?username=John&email=example@example.com
//7. Пагинация с использованием параметров page и page_size (например, 2-я страница с 5 пользователями на странице):
//GET /users?page=2&page_size=5
//8. Комбинация фильтрации и пагинации:
//GET /users?status=active&created_after=2024-01-01&page=1&page_size=10

//Получить пользователей, созданных после 2024-01-01:
//GET /users?created_after=2024-01-01
//Получить пользователей, созданных до 2024-12-31:
//GET /users?created_before=2024-12-31
//Получить пользователей, созданных в 2024 году:
//GET /users?created_after=2024-01-01&created_before=2024-12-31

func GetUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	filter := bson.M{}

	// Фильтрация по имени пользователя (например, ?username=John)
	if username := queryParams.Get("username"); username != "" {
		filter["username"] = bson.M{"$regex": username, "$options": "i"} // Нерегистрозависимый поиск
	}

	// Фильтрация по email (например, ?email=example@example.com)
	if email := queryParams.Get("email"); email != "" {
		filter["email"] = bson.M{"$regex": email, "$options": "i"} // Нерегистрозависимый поиск
	}

	// Фильтрация по дате создания (диапазон)
	createdAtFilter := bson.M{}
	if createdAfter := queryParams.Get("created_after"); createdAfter != "" {
		createdAfterTime, err := time.Parse("2006-01-02", createdAfter)
		if err != nil {
			http.Error(w, "Invalid created_after date format", http.StatusBadRequest)
			return
		}
		createdAtFilter["$gte"] = createdAfterTime
	}
	if createdBefore := queryParams.Get("created_before"); createdBefore != "" {
		createdBeforeTime, err := time.Parse("2006-01-02", createdBefore)
		if err != nil {
			http.Error(w, "Invalid created_before date format", http.StatusBadRequest)
			return
		}
		createdAtFilter["$lte"] = createdBeforeTime
	}
	if len(createdAtFilter) > 0 {
		filter["created_at"] = createdAtFilter
	}

	// Фильтрация по статусу аккаунта (например, ?status=active)
	if status := queryParams.Get("status"); status != "" {
		filter["status"] = bson.M{"$regex": status, "$options": "i"} // Нерегистрозависимый поиск
	}

	usersCollection := db.GetCollection("users")

	// Пагинация
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

//GET /users/{id}

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

/*
PUT /users/{id}
Content-Type: application/json

{
"username": "new_username",
"email": "new_email@example.com",
"permissions": "admin"
}
*/

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// PATCH /users/5fcbf22b923e992dcebf6f1a
// В теле запроса клиент может отправить те поля пользователя, которые он хочет обновить (например, username, email, permissions, payments, jobs).
// Важно, что только те поля, которые не пустые, будут включены в обновление.
/*
{
    "username": "john_doe_updated",
    "email": "john.doe.updated@example.com",
    "permissions": "admin",
    "payments": [
        {
            "payment_id": "12345",
            "amount": 100,
            "currency": "USD"
        }
    ],
    "jobs": [
        "5fcbf22b923e992dcebf6f1b"
    ]
}

*/

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
	if len(patchUser.Payments) > 0 {
		updateFields["payments"] = patchUser.Payments
	}
	if len(patchUser.Jobs) > 0 {
		updateFields["jobs"] = patchUser.Jobs
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

/*
	{
	    "username": "john_doe",
	    "email": "john.doe@example.com",
	    "password_hash": "hashed_password_value",
	    "permissions": "admin"
	}

Ответ:

	{
	    "id": "5fcbf22b923e992dcebf6f1a",
	    "username": "john_doe",
	    "email": "john.doe@example.com",
	    "permissions": "admin",
	    "created_at": "2024-12-08T00:00:00Z",
	    "updated_at": "2024-12-08T00:00:00Z"
	}
*/
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

	if newUser.Permissions == "" {
		newUser.Permissions = "user" // по умолчанию
	}

	newUser.ID = primitive.NewObjectID()
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()

	usersCollection := db.GetCollection("users")

	_, err := usersCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		http.Error(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// DELETE /users/{userID}
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

// GET /users/{userID}/jobs

//Ответ - массив объектов Job

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
	err = UpdateJobsStatus(jobsCollection)
	if err != nil {
		http.Error(w, "Error while updating jobs", http.StatusInternalServerError)
		return
	}
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

/*
POST /users/{userID}/jobs
{
  "title": "New Translation Job",
  "status": "pending",
  "source_language": "en",
  "file_format": "pdf",
  "description": "Translate a document from English to Spanish.",
  "input_file": "input_file_path.pdf",
  "output_file": "output_file_path.pdf"
}
*/

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

	const FinishDatetime = 30

	job.ID = primitive.NewObjectID()
	job.UserID = id
	job.CreatedAt = time.Now()
	job.UpdatedAt = job.CreatedAt
	job.EstimatedFinishDatetime = time.Now().Add(FinishDatetime * time.Second)

	if job.Title == "" || job.Status == "" || job.SourceLanguage == "" || job.FileFormat == "" || job.Description == "" || job.InputFile == "" || job.OutputFile == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

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

	var user models.User
	err = usersCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	if user.Jobs == nil {
		user.Jobs = []primitive.ObjectID{job.ID}
	} else {
		user.Jobs = append(user.Jobs, job.ID)
	}

	_, err = usersCollection.UpdateOne(context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"jobs": user.Jobs}},
	)
	if err != nil {
		http.Error(w, "Error updating user jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}

// DELETE /users/{id}/jobs/{job_id}
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

/*
Запрос:
POST /users/60d09c875d3b3c6b8d85a681/payments
Content-Type: application/json

Тело запроса:
{
	"price": "100.00",
	"payment_method": "credit_card",
	"payment_status": "completed",
	"job_id": "60d09c875d3b3c6b8d85a683"
}

Ответ:
{
	"id": "60d09c875d3b3c6b8d85a685",
	"price": "100.00",
	"payment_method": "credit_card",
	"payment_status": "completed",
	"created_at": "2024-12-08T12:00:00Z",
	"updated_at": "2024-12-08T12:00:00Z",
	"job_id": "60d09c875d3b3c6b8d85a683"
}
*/

func AddPayment(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	var payment models.Payment
	if err := render.DecodeJSON(r.Body, &payment); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()

	filter := bson.M{"_id": objectID}

	update := bson.M{
		"$push": bson.M{
			"payments": payment,
		},
	}

	_, err = db.GetCollection("users").UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error adding payment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payment)
}

// DELETE /users/60d09c875d3b3c6b8d85a681/payments/60d09c875d3b3c6b8d85a685
func DeletePayment(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	paymentID := chi.URLParam(r, "payment_id")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	paymentObjectID, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		http.Error(w, "Invalid payment ID format", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$pull": bson.M{
			"payments": bson.M{"_id": paymentObjectID},
		},
	}

	_, err = db.GetCollection("users").UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error deleting payment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateJobsStatus(jobsCollection *mongo.Collection) error {
	filter := bson.M{
		"status": bson.M{"$ne": "completed"},
	}

	cursor, err := jobsCollection.Find(context.Background(), filter)
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	currentTime := time.Now()

	// Пробегаем по задачам.
	for cursor.Next(context.Background()) {
		var job struct {
			ID                      interface{} `bson:"_id"`
			EstimatedFinishDatetime time.Time   `bson:"estimated_finish_datetime"`
			Status                  string      `bson:"status"`
		}

		if err := cursor.Decode(&job); err != nil {
			log.Printf("Error decoding job: %v", err)
			continue
		}

		// Если задача выполнена, обновляем её статус.
		if job.EstimatedFinishDatetime.Before(currentTime) && job.Status != "completed" {
			update := bson.M{
				"$set": bson.M{
					"status":     "completed",
					"updated_at": currentTime,
				},
			}

			_, err := jobsCollection.UpdateOne(context.Background(), bson.M{"_id": job.ID}, update)
			if err != nil {
				log.Printf("Error updating job %v: %v", job.ID, err)
				continue
			}
			log.Printf("Updated job %v to completed", job.ID)
		}
	}
	if err := cursor.Err(); err != nil {
		return err
	}
	return nil
}

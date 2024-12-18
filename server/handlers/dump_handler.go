package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/moevm/nosql2h24-transcribtion/db"
	"github.com/moevm/nosql2h24-transcribtion/models"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

// Структура для экспорта/импорта данных
type SystemData struct {
	Users   []models.User   `json:"users"`
	Servers []models.Server `json:"servers"`
	Jobs    []models.Job    `json:"jobs"`
}

// Экспорт данных
func ExportData(w http.ResponseWriter, r *http.Request) {
	usersCollection := db.GetCollection("users")
	serversCollection := db.GetCollection("servers")
	jobsCollection := db.GetCollection("jobs")

	var systemData SystemData

	// Получение данных из коллекции users
	usersCursor, err := usersCollection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Error exporting users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer usersCursor.Close(context.Background())

	var users []models.User
	if err := usersCursor.All(context.Background(), &users); err != nil {
		http.Error(w, "Error decoding users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	systemData.Users = users

	// Получение данных из коллекции servers
	serversCursor, err := serversCollection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Error exporting servers: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer serversCursor.Close(context.Background())

	var servers []models.Server
	if err := serversCursor.All(context.Background(), &servers); err != nil {
		http.Error(w, "Error decoding servers: "+err.Error(), http.StatusInternalServerError)
		return
	}
	systemData.Servers = servers

	// Получение данных из коллекции jobs
	jobsCursor, err := jobsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		http.Error(w, "Error exporting jobs: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer jobsCursor.Close(context.Background())

	var jobs []models.Job
	if err := jobsCursor.All(context.Background(), &jobs); err != nil {
		http.Error(w, "Error decoding jobs: "+err.Error(), http.StatusInternalServerError)
		return
	}
	systemData.Jobs = jobs

	// Отправка данных в формате JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(systemData); err != nil {
		http.Error(w, "Error encoding data: "+err.Error(), http.StatusInternalServerError)
	}
}

func ImportData(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Чтение тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Парсинг данных
	var systemData SystemData
	if err := json.Unmarshal(body, &systemData); err != nil {
		http.Error(w, "Error parsing JSON data: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Коллекции для работы с MongoDB
	usersCollection := db.GetCollection("users")
	serversCollection := db.GetCollection("servers")
	jobsCollection := db.GetCollection("jobs")

	// Очистка всех коллекций
	clearCollection := func(collection *mongo.Collection, collectionName string) error {
		_, err := collection.DeleteMany(ctx, bson.M{})
		if err != nil {
			return fmt.Errorf("Error clearing %s: %v", collectionName, err)
		}
		return nil
	}

	// Очистка коллекций перед импортом
	if err := clearCollection(usersCollection, "users"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := clearCollection(serversCollection, "servers"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := clearCollection(jobsCollection, "jobs"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Функция для вставки данных в коллекцию
	insertMany := func(collection *mongo.Collection, data interface{}, collectionName string) error {
		// Преобразование data в []interface{}
		var dataSlice []interface{}
		switch v := data.(type) {
		case []models.User:
			for _, user := range v {
				dataSlice = append(dataSlice, user)
			}
		case []models.Server:
			for _, server := range v {
				dataSlice = append(dataSlice, server)
			}
		case []models.Job:
			for _, job := range v {
				dataSlice = append(dataSlice, job)
			}
		default:
			return fmt.Errorf("unsupported data type for %s", collectionName)
		}

		_, err := collection.InsertMany(ctx, dataSlice)
		if err != nil {
			return fmt.Errorf("Error importing %s: %v", collectionName, err)
		}
		return nil
	}

	// Импорт данных в коллекции
	if err := insertMany(usersCollection, systemData.Users, "users"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := insertMany(serversCollection, systemData.Servers, "servers"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := insertMany(jobsCollection, systemData.Jobs, "jobs"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Ответ об успешном импорте данных
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data imported successfully"))
}

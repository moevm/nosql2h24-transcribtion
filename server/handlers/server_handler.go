package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/moevm/nosql2h24-transcribtion/db"
	"github.com/moevm/nosql2h24-transcribtion/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
	"time"
)

// GetServers - обработчик для получения списка серверов с фильтрацией
// @Description Возвращает список серверов с поддержкой фильтрации по CPU, GPU, RAM, статусу
// @Param status "Фильтр по статусу сервера"
// @Param cpu"Фильтр по CPU"
// @Param gpu "Фильтр по GPU"
// @Param ram query int  "Фильтр по МИНИМАЛЬНОМУ объему RAM в ГБ"
func GetServers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	status := queryParams.Get("status")
	cpu := queryParams.Get("cpu")
	gpu := queryParams.Get("gpu")
	ramStr := queryParams.Get("ram")

	var ram int
	var err error

	if ramStr != "" {
		ram, err = strconv.Atoi(ramStr)
		if err != nil {
			http.Error(w, "Invalid RAM value", http.StatusBadRequest)
			return
		}
	}

	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}
	if cpu != "" {
		filter["cpu_info"] = bson.M{"$regex": cpu, "$options": "i"}
	}
	if gpu != "" {
		filter["gpu_info"] = bson.M{"$regex": gpu, "$options": "i"}
	}
	if ram > 0 {
		filter["ram_size_gb"] = bson.M{"$gte": ram}
	}

	collection := db.GetCollection("servers")
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, "Error fetching servers", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var servers []models.Server
	if err = cursor.All(context.Background(), &servers); err != nil {
		http.Error(w, "Error decoding servers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(servers)
}

func GetServerByID(w http.ResponseWriter, r *http.Request) {
	serverID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(serverID)
	if err != nil {
		http.Error(w, "Invalid server ID", http.StatusBadRequest)
		return
	}
	serversCollection := db.GetCollection("servers")

	var server models.Server
	err = serversCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&server)
	if err != nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(server)
}

func CreateServer(w http.ResponseWriter, r *http.Request) {
	var newServer models.Server
	if err := render.DecodeJSON(r.Body, &newServer); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if newServer.Hostname == "" || newServer.Address == "" { //Уточнить, какие поля обязательные
		http.Error(w, "Hostname and Address are required", http.StatusBadRequest)
		return
	}

	newServer.CreatedAt = time.Now()
	newServer.UpdatedAt = time.Now()

	newServer.ID = primitive.NewObjectID()
	serversCollection := db.GetCollection("servers")

	_, err := serversCollection.InsertOne(context.Background(), newServer)
	if err != nil {
		http.Error(w, "Error saving server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newServer)
}

func UpdateServer(w http.ResponseWriter, r *http.Request) {
	serverID := chi.URLParam(r, "id")
	objectID, err := primitive.ObjectIDFromHex(serverID)
	if err != nil {
		http.Error(w, "Invalid server ID", http.StatusBadRequest)
		return
	}

	var updatedServer models.Server
	if err := render.DecodeJSON(r.Body, &updatedServer); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if updatedServer.Hostname == "" || updatedServer.Address == "" {
		http.Error(w, "Hostname and Address are required", http.StatusBadRequest)
		return
	}
	updatedServer.UpdatedAt = time.Now()
	serversCollection := db.GetCollection("servers")

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": updatedServer}

	_, err = serversCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error updating server", http.StatusInternalServerError)
		return
	}

	updatedServer.ID = objectID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedServer)
}

func PatchServer(w http.ResponseWriter, r *http.Request) {
	serverID := chi.URLParam(r, "id")
	objectID, err := primitive.ObjectIDFromHex(serverID)
	if err != nil {
		http.Error(w, "Invalid server ID", http.StatusBadRequest)
		return
	}

	var patchData models.Server
	if err := render.DecodeJSON(r.Body, &patchData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	patchData.UpdatedAt = time.Now()
	serversCollection := db.GetCollection("servers")

	update := bson.M{}
	if patchData.Hostname != "" {
		update["hostname"] = patchData.Hostname
	}
	if patchData.Address != "" {
		update["address"] = patchData.Address
	}
	if patchData.Description != "" {
		update["description"] = patchData.Description
	}
	if patchData.Status != "" {
		update["status"] = patchData.Status
	}
	if patchData.CPUInfo != "" {
		update["cpu_info"] = patchData.CPUInfo
	}
	if patchData.GPUInfo != "" {
		update["gpu_info"] = patchData.GPUInfo
	}
	if patchData.RAMSizeGB != 0 {
		update["ram_size_gb"] = patchData.RAMSizeGB
	}

	if len(update) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": objectID}
	updateQuery := bson.M{"$set": update}

	_, err = serversCollection.UpdateOne(context.Background(), filter, updateQuery)
	if err != nil {
		http.Error(w, "Error updating server", http.StatusInternalServerError)
		return
	}

	patchData.ID = objectID // Устанавливаем ID сервера
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patchData)
}

func DeleteServer(w http.ResponseWriter, r *http.Request) {
	serverID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(serverID)
	if err != nil {
		http.Error(w, "Invalid server ID", http.StatusBadRequest)
		return
	}

	serversCollection := db.GetCollection("servers")
	result, err := serversCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, "Error deleting server", http.StatusInternalServerError)
		return
	}
	if result.DeletedCount == 0 {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Server deleted successfully"})
}

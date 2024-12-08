package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/moevm/nosql2h24-transcribtion/db"
	"github.com/moevm/nosql2h24-transcribtion/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

// GET /servers?status=active&cpu=Intel&gpu=NVIDIA&ram=16
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

// GET /servers/{id}
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

/*
PUT /servers/607f1f77bcf86cd799439011

	{
	  "hostname": "new-server-name",
	  "address": "192.168.1.10",
	  "status": "inactive",
	  "cpu_info": "Intel Xeon E5",
	  "ram_size_gb": 64
	}
*/
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

/*
PATCH /servers/607f1f77bcf86cd799439011

	{
	  "hostname": "new-server-name",
	  "address": "192.168.1.10",
	  "description": "Updated server description",
	  "status": "inactive",
	  "cpu_info": "Intel Xeon E5",
	  "gpu_info": "NVIDIA Tesla",
	  "ram_size_gb": 64
	}
*/
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

	patchData.ID = objectID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patchData)
}

// DELETE /servers/{id}
func DeleteServer(w http.ResponseWriter, r *http.Request) {
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
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "Server not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching server details: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if len(server.CurrentJobs) > 0 || len(server.CompletedJobs) > 0 {
		http.Error(w, "Cannot delete server with associated jobs", http.StatusConflict)
		return
	}

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

// Get /servers/{id}/currentJobs
func GetServerCurrentJobs(w http.ResponseWriter, r *http.Request) {
	serverID := chi.URLParam(r, "id")
	serversCollection := db.GetCollection("servers")

	serverIDObj, err := primitive.ObjectIDFromHex(serverID)
	if err != nil {
		http.Error(w, "Invalid server ID", http.StatusBadRequest)
		return
	}

	var server models.Server
	err = serversCollection.FindOne(context.Background(), bson.M{"_id": serverIDObj}).Decode(&server)
	if err != nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}

	jobsCollection := db.GetCollection("jobs")
	cursor, err := jobsCollection.Find(context.Background(), bson.M{"_id": bson.M{"$in": server.CurrentJobs}})
	if err != nil {
		http.Error(w, "Error fetching current jobs", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var jobs []models.Job
	err = cursor.All(context.Background(), &jobs)
	if err != nil {
		http.Error(w, "Error processing current jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

// Get /servers/{id}/completedJobs
func GetServerCompletedJobs(w http.ResponseWriter, r *http.Request) {
	serverID := chi.URLParam(r, "id")
	serversCollection := db.GetCollection("servers")

	serverIDObj, err := primitive.ObjectIDFromHex(serverID)
	if err != nil {
		http.Error(w, "Invalid server ID", http.StatusBadRequest)
		return
	}

	var server models.Server
	err = serversCollection.FindOne(context.Background(), bson.M{"_id": serverIDObj}).Decode(&server)
	if err != nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}

	jobsCollection := db.GetCollection("jobs")
	cursor, err := jobsCollection.Find(context.Background(), bson.M{"_id": bson.M{"$in": server.CompletedJobs}})
	if err != nil {
		http.Error(w, "Error fetching completed jobs", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var jobs []models.Job
	err = cursor.All(context.Background(), &jobs)
	if err != nil {
		http.Error(w, "Error processing completed jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func AddJobToServer(w http.ResponseWriter, r *http.Request) {
	serverID := chi.URLParam(r, "id")
	jobID := chi.URLParam(r, "job_id")
	serversCollection := db.GetCollection("servers")
	jobsCollection := db.GetCollection("jobs")

	serverIDObj, err := primitive.ObjectIDFromHex(serverID)
	if err != nil {
		http.Error(w, "Invalid server ID", http.StatusBadRequest)
		return
	}

	jobIDObj, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	var server models.Server
	err = serversCollection.FindOne(context.Background(), bson.M{"_id": serverIDObj}).Decode(&server)
	if err != nil {
		http.Error(w, "Server not found", http.StatusNotFound)
		return
	}

	var job models.Job
	err = jobsCollection.FindOne(context.Background(), bson.M{"_id": jobIDObj}).Decode(&job)
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	_, err = serversCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": serverIDObj},
		bson.M{"$push": bson.M{"current_jobs": jobIDObj}},
	)
	if err != nil {
		http.Error(w, "Error adding job to server", http.StatusInternalServerError)
		return
	}

	_, err = jobsCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": jobIDObj},
		bson.M{"$set": bson.M{"host_id": serverIDObj}},
	)
	if err != nil {
		http.Error(w, "Error updating job with server info", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Job successfully added to server"})
}

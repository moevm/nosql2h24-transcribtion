package handlers

import (
	"context"
	"encoding/json"
	"github.com/moevm/nosql2h24-transcribtion/db"
	"github.com/moevm/nosql2h24-transcribtion/models"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
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

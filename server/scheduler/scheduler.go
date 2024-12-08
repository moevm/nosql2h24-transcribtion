package schedul

import (
	"context"
	"errors"
	"github.com/moevm/nosql2h24-transcribtion/models"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetServers(serversCollection *mongo.Collection) ([]models.Server, error) {
	var servers []models.Server
	cursor, err := serversCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.Background(), &servers); err != nil {
		return nil, err
	}
	if len(servers) == 0 {
		return nil, errors.New("no servers available")
	}
	return servers, nil
}

func SelectServerWithMinJobs(servers []models.Server) (models.Server, error) {
	if len(servers) == 0 {
		return models.Server{}, errors.New("no servers provided")
	}

	var selectedServers []models.Server
	minJobs := int(^uint(0) >> 1)
	for _, server := range servers {
		if len(server.CurrentJobs) < minJobs {
			minJobs = len(server.CurrentJobs)
			selectedServers = []models.Server{server}
		} else if len(server.CurrentJobs) == minJobs {
			selectedServers = append(selectedServers, server)
		}
	}

	rand.Seed(time.Now().UnixNano())
	return selectedServers[rand.Intn(len(selectedServers))], nil
}

func AddJobToServer(serversCollection *mongo.Collection, serverID primitive.ObjectID, jobID primitive.ObjectID) error {
	_, err := serversCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": serverID},
		bson.M{"$push": bson.M{"current_jobs": jobID}},
	)
	return err
}

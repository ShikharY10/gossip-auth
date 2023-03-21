package handlers

import (
	"context"
	"errors"

	"github.com/ShikharY10/gbAUTH/cmd/admin"
	config "github.com/ShikharY10/gbAUTH/cmd/configs"
	"github.com/ShikharY10/gbAUTH/cmd/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DataBase struct {
	userCollection      *mongo.Collection
	frequencyCollection *mongo.Collection
	payloadsCollection  *mongo.Collection
}

func InitializeDataBase(mongodb *config.MongoDB, logger *admin.Logger) *DataBase {
	return &DataBase{
		userCollection:      mongodb.User,
		frequencyCollection: mongodb.Frequnecy,
		payloadsCollection:  mongodb.Payloads,
	}
}

func (db *DataBase) IsUsernameAwailable(username string) error {
	cursor, err := db.frequencyCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return err
	}
	var users []models.FrequencyTable
	for cursor.Next(context.TODO()) {
		var elem models.FrequencyTable
		err := cursor.Decode(&elem)
		if err != nil {
			return err
		}
		users = append(users, elem)
	}
	for _, user := range users {
		if user.Username == username {
			return errors.New("username already present")
		}
	}
	return nil
}

package handlers

import (
	"context"
	"errors"
	"time"

	"github.com/ShikharY10/gbAUTH/cmd/admin"
	config "github.com/ShikharY10/gbAUTH/cmd/configs"
	"github.com/ShikharY10/gbAUTH/cmd/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (db *DataBase) IsEmailAvailable(email string) error {
	opts := options.Find().SetProjection(bson.D{{Key: "_id", Value: 1}})
	cursor, err := db.userCollection.Find(
		context.TODO(),
		bson.M{"email": email},
		opts,
	)
	if err != nil {
		return err
	}
	var users []models.User
	for cursor.Next(context.TODO()) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			return err
		}
		users = append(users, user)
	}
	if len(users) == 0 {
		return nil
	}
	return errors.New("email already exist")

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

func (db *DataBase) AddUserPayloadsField() (*primitive.ObjectID, error) {
	b := bson.M{
		"msg": bson.M{},
	}
	res, err := db.payloadsCollection.InsertOne(context.TODO(), b)
	if err != nil {
		return nil, err
	}
	_id := res.InsertedID.(primitive.ObjectID)
	return &_id, nil
}

func (db *DataBase) CreateNewUser(user models.User) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFunc()
	_, err := db.userCollection.InsertOne(ctx, user)
	return err
}

func (db *DataBase) InsetUserInFrequencyTable(id primitive.ObjectID, username string) error {
	fTable := models.FrequencyTable{
		Id:        id,
		Username:  username,
		Frequency: 0,
	}
	_, err := db.frequencyCollection.InsertOne(
		context.TODO(),
		fTable,
	)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (db *DataBase) GetUserData(filter bson.M, findOptions *options.FindOneOptions) (*models.User, error) {
	cursor := db.userCollection.FindOne(context.TODO(), filter, findOptions)

	if cursor.Err() != nil {
		return nil, cursor.Err()
	}

	var user models.User
	err := cursor.Decode(&user)
	if err != nil {
		return nil, err
	} else {

		return &user, nil
	}
}

func (db *DataBase) GetUsersData(filter bson.M, findOption *options.FindOptions) (*[]models.User, error) {
	cursor, err := db.userCollection.Find(context.TODO(), filter, findOption)
	if err != nil {
		return nil, err
	}

	var users []models.User
	err = cursor.All(context.TODO(), &users)
	if err != nil {
		return nil, err
	}

	if len(users) > 0 {
		return &users, nil
	}
	return nil, errors.New("no document found")
}

func (db *DataBase) UpdateLogoutStatus(username string, status bool) error {
	result, err := db.userCollection.UpdateOne(
		context.TODO(),
		bson.M{"username": username},
		bson.M{"$set": bson.M{"logout": status}},
	)
	if err != nil {
		return err
	}

	if result.ModifiedCount > int64(0) {
		return nil
	}
	return err
}

func (db *DataBase) GetUserEmail(id string) (string, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}
	otps := options.FindOne().SetProjection(bson.D{{Key: "email", Value: 1}})
	user, err := db.GetUserData(bson.M{"_id": _id}, otps)
	if err != nil {
		return "", err
	} else {
		return user.Email, nil
	}
}

func (db *DataBase) UpdateUserName(id string, name string) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	} else {
		result, err := db.userCollection.UpdateOne(
			context.TODO(),
			bson.M{"_id": _id},
			bson.M{"$set": bson.M{"name": name}},
		)
		if err != nil {
			return err
		}
		if result.ModifiedCount > int64(0) {
			return nil
		}
		return errors.New("something went wrong")
	}
}

func (db *DataBase) UpdateUserAvatar(id string, avatar models.Avatar) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	} else {
		result, err := db.userCollection.UpdateOne(
			context.TODO(),
			bson.M{"_id": _id},
			bson.M{"$set": bson.M{"avatar": avatar}},
		)
		if err != nil {
			return err
		}
		if result.ModifiedCount > int64(0) {
			return nil
		}
		return errors.New("something went wrong")
	}
}

func (db *DataBase) UpdateUserDetail(id string, key string, value any) error {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	} else {
		result, err := db.userCollection.UpdateOne(
			context.TODO(),
			bson.M{"_id": _id},
			bson.M{"$set": bson.M{key: value}},
		)
		if err != nil {
			return err
		}
		if result.ModifiedCount > int64(0) {
			return nil
		}
		return errors.New("something went wrong")
	}
}

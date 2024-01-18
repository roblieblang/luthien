package dao

import (
	"context"

	"github.com/roblieblang/luthien-core-server/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

type UserDAO struct {
	collection *mongo.Collection
}

func NewUserDAO(client *mongo.Client, dbName, collectionName string) *UserDAO {
	collection := client.Database(dbName).Collection(collectionName)
	return &UserDAO{collection: collection}
}

func (dao *UserDAO) CreateUser(user *models.User) error {
	_, err := dao.collection.InsertOne(context.TODO(), user)
	return err
}

func (dao *UserDAO) GetUser(id int) (*models.User, error) {
	var user models.User
	filter := bson.M{"id": id}
	err := dao.collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

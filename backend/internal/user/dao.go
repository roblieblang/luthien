package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DAO struct {
	collection *mongo.Collection
}

func NewDAO(client *mongo.Client, dbName, collectionName string) *DAO {
	collection := client.Database(dbName).Collection(collectionName)
	return &DAO{collection: collection}
}

func (dao *DAO) CreateUser(user *User) error {
	_, err := dao.collection.InsertOne(context.TODO(), user)
	return err
}

func (dao *DAO) GetUser(id string) (*User, error) {
	var user User
	objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }
	filter := bson.M{"_id": objID}
	err = dao.collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *DAO) GetAllUsers() ([]User,error) {
	var users []User
	cursor, err := dao.collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
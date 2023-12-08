package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Username  string    `json:"username" bson:"username"`
    Email     string    `json:"email" bson:"email"`
    Password  string    `json:"-" bson:"password"`
    FirstName string    `json:"firstName" bson:"firstName,omitempty"`
    LastName  string    `json:"lastName" bson:"lastName,omitempty"`
    CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}
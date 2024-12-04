package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/golang-jwt/jwt/v5"
	"time"
)


type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Password []byte             `json:"-" bson:"password"`
	RoleID   primitive.ObjectID `json:"role_id" bson:"role_id"` 
}


type Task struct {
    ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
    Name      string             `json:"name" bson:"name"`
    Description string           `json:"description" bson:"description"`
    Status    bool             `json:"status" bson:"status"`
    CreatedAt time.Time          `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}


type Permission struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`         
	Description string             `json:"description" bson:"description"`
}


type Role struct {
	ID          primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string               `json:"name" bson:"name"`       
	Permissions []primitive.ObjectID `json:"permissions" bson:"permissions"` 
}


type CustomClaims struct {
    Role   string             `json:"role"`
    jwt.RegisteredClaims
}

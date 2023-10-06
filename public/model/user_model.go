package model

import (
	"time"

	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/proto/pb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PublicUserModel struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type CreateUserModel struct {
	Email    string `json:"email" bson:"email"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

func (c *CreateUserModel) ToCreateUserRequest() *pb.CreateUserRequest {
	return &pb.CreateUserRequest{
		Email:    c.Email,
		Username: c.Username,
		Password: c.Password,
	}
}

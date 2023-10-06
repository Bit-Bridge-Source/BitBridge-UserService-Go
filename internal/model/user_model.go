package model

import (
	"time"

	"github.com/Bit-Bridge-Source/BitBridge-UserService-Go/public/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PrivateUserModel struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email"`
	Username  string             `json:"username" bson:"username"`
	Hash      string             `json:"password" bson:"password"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func (privateUserModel *PrivateUserModel) ToPublicUserModel() *model.PublicUserModel {
	return &model.PublicUserModel{
		ID:        privateUserModel.ID,
		Username:  privateUserModel.Username,
		CreatedAt: privateUserModel.CreatedAt,
	}
}

package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CartProduct struct {
	ProductID primitive.ObjectID `json:"productId" bson:"productId"`
	Quantity  int                `json:"quantity" bson:"quantity" validate:"required,gt=0"`
}

type Cart struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID   primitive.ObjectID `json:"userId" bson:"userId"`
	Products []CartProduct      `json:"products" bson:"products"`
}

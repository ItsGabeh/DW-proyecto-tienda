package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name" validate:"required"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price" validate:"required,gt=0"`
	Stock       int                `json:"stock" bson:"stock" validate:"required,gte=0"`
}

package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
    ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Title  string             `bson:"title" json:"title"`
    ISBN   string             `bson:"isbn" json:"isbn"`
    Author string             `bson:"author" json:"author"`
}
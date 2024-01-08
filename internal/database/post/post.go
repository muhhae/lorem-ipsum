package post

import (
	"context"
	"time"

	"github.com/muhhae/lorem-ipsum/internal/database/connection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	AuthorID  primitive.ObjectID   `json:"author" bson:"author,required"`
	Content   string               `json:"content" bson:"content,required"`
	ImageIDs  []primitive.ObjectID `json:"image" bson:"image,required"`
	CreatedAt time.Time            `json:"createdAt" bson:"createdAt,required"`
}

func (p *Post) Save() (primitive.ObjectID, error) {
	posts := connection.GetDB().Posts
	if p.ID == primitive.NilObjectID {
		_, err := posts.InsertOne(context.Background(), p)
		return primitive.NilObjectID, err
	}
	count, err := posts.CountDocuments(context.Background(), bson.M{"_id": p.ID})
	if err != nil {
		return primitive.NilObjectID, err

	}
	if count > 0 {
		res := posts.FindOneAndReplace(context.Background(), bson.M{"_id": p.ID}, p)
		return primitive.NilObjectID, res.Err()
	}
	id, err := posts.InsertOne(context.Background(), p)
	return id.InsertedID.(primitive.ObjectID), err
}

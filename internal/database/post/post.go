package post

import (
	"context"
	"time"

	"github.com/muhhae/lorem-ipsum/internal/database/connection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	AuthorID  primitive.ObjectID `json:"author" bson:"author,required"`
	Content   string             `json:"content" bson:"content,required"`
	Image     string             `json:"image" bson:"image,required"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt,required"`
}

func (p *Post) Save() error {
	posts := connection.GetDB().Posts
	if p.ID == primitive.NilObjectID {
		_, err := posts.InsertOne(context.Background(), p)
		return err
	}
	count, err := posts.CountDocuments(context.Background(), bson.M{"_id": p.ID})
	if err != nil {
		return err
	}
	if count > 0 {
		res := posts.FindOneAndReplace(context.Background(), bson.M{"_id": p.ID}, p)
		return res.Err()
	}
	_, err = posts.InsertOne(context.Background(), p)
	return err
}

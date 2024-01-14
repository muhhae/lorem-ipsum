package comment

import (
	"context"
	"time"

	"github.com/muhhae/lorem-ipsum/internal/database/connection"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PostID    primitive.ObjectID `json:"post_id" bson:"post_id,required"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id,required"`
	Content   string             `json:"content" bson:"content,required"`
	Parent    primitive.ObjectID `json:"parent" bson:"parent,required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,required"`
}

func (c *Comment) Save() (primitive.ObjectID, error) {
	col := connection.GetDB().Comments
	if c.ID != primitive.NilObjectID {
		res := col.FindOneAndReplace(context.Background(), primitive.M{"_id": c.ID}, c)
		if res.Err() != nil {
			return primitive.NilObjectID, res.Err()
		}
		return c.ID, nil
	}
	c

	res, err := col.InsertOne(context.Background(), c)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

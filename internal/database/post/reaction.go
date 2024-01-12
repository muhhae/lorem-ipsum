package post

import (
	"context"

	"github.com/muhhae/lorem-ipsum/internal/database/connection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Reaction struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PostID primitive.ObjectID `bson:"post_id" json:"post_id"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	Value  int                `bson:"value" json:"value"`
}

func (r *Reaction) Save() (primitive.ObjectID, error) {
	col := connection.GetDB().Reactions
	res := col.FindOne(context.Background(), bson.M{"post_id": r.PostID, "user_id": r.UserID})
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			res, err := col.InsertOne(context.Background(), r)
			if err != nil {
				return primitive.NilObjectID, err
			}
			return res.InsertedID.(primitive.ObjectID), nil
		}
		return primitive.NilObjectID, res.Err()
	}
	currentReaction := Reaction{}
	err := res.Decode(&currentReaction)
	if err != nil {
		return primitive.NilObjectID, err
	}
	res = col.FindOneAndReplace(context.Background(), bson.M{"_id": currentReaction.ID}, r)
	if res.Err() != nil {
		return primitive.NilObjectID, res.Err()
	}
	return currentReaction.ID, nil
}

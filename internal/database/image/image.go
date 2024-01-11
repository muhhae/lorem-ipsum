package imagemodel

import (
	"context"

	"github.com/muhhae/lorem-ipsum/internal/database/connection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Image struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Owner primitive.ObjectID `json:"owner" bson:"owner,required"`
	Data  []byte             `json:"data" bson:"data,required"`
}

func (img *Image) Save() (primitive.ObjectID, error) {
	imgs := connection.GetDB().Images
	if img.ID == primitive.NilObjectID {
		inserted, err := imgs.InsertOne(context.Background(), img)
		return inserted.InsertedID.(primitive.ObjectID), err
	}
	count, err := imgs.CountDocuments(context.Background(), bson.M{"_id": img.ID})
	if err != nil {
		return primitive.NilObjectID, err
	}
	if count > 0 {
		res := imgs.FindOneAndReplace(context.Background(), bson.M{"_id": img.ID}, img)
		return primitive.NilObjectID, res.Err()
	}
	imgID, err := imgs.InsertOne(context.Background(), img)
	return imgID.InsertedID.(primitive.ObjectID), err
}

func FindOne(filter bson.M) (*Image, error) {
	imgs := connection.GetDB().Images
	img := Image{}
	res := imgs.FindOne(context.Background(), filter)
	if res.Err() != nil {
		return nil, res.Err()
	}
	err := res.Decode(&img)
	if err != nil {
		return nil, err
	}
	return &img, nil
}

package image

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

func (img *Image) Save() error {
	imgs := connection.GetDB().Images
	if img.ID == primitive.NilObjectID {
		_, err := imgs.InsertOne(context.Background(), img)
		return err
	}
	count, err := imgs.CountDocuments(context.Background(), bson.M{"_id": img.ID})
	if err != nil {
		return err
	}
	if count > 0 {
		res := imgs.FindOneAndReplace(context.Background(), bson.M{"_id": img.ID}, img)
		return res.Err()
	}
	_, err = imgs.InsertOne(context.Background(), img)
	return err
}

func FindOne(filter interface{}) (*Image, error) {
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

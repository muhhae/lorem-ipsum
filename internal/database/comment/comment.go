package comment

import (
	"context"
	"time"

	"github.com/muhhae/lorem-ipsum/internal/database/connection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	res, err := col.InsertOne(context.Background(), c)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func FindAll(filter bson.M) ([]Comment, error) {
	col := connection.GetDB().Comments
	cur, err := col.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var comments []Comment
	for cur.Next(context.Background()) {
		var comment Comment
		if err := cur.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

const (
	maxComments = 5
)

func RetrieveDefault(postID primitive.ObjectID, iteration int) ([]Comment, error) {
	var comments []Comment
	col := connection.GetDB().Comments
	skip := int64(iteration * maxComments)
	limit := int64(maxComments)
	filter := primitive.M{
		"post_id": postID,
		"parent":  primitive.NilObjectID,
	}
	cursor, err := col.Find(context.Background(), filter, &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
		Sort:  bson.M{"created_at": -1},
	})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.Background()) {
		var comment Comment
		if err := cursor.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func RetrieveAll(postID primitive.ObjectID) ([]Comment, error) {
	var comments []Comment
	col := connection.GetDB().Comments
	filter := primitive.M{
		"post_id": postID,
		"parent":  primitive.NilObjectID,
	}
	cursor, err := col.Find(context.Background(), filter, &options.FindOptions{
		Sort: bson.M{"created_at": -1},
	})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.Background()) {
		var comment Comment
		if err := cursor.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func RetrieveUserComments(userID primitive.ObjectID, postID primitive.ObjectID) ([]Comment, error) {
	var comments []Comment
	col := connection.GetDB().Comments
	filter := primitive.M{
		"user_id": userID,
		"post_id": postID,
		"parent":  primitive.NilObjectID,
	}
	cursor, err := col.Find(context.Background(), filter, &options.FindOptions{
		Sort: bson.M{"created_at": -1},
	})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.Background()) {
		var comment Comment
		if err := cursor.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func CommentCount(postID primitive.ObjectID) (int64, error) {
	col := connection.GetDB().Comments
	filter := primitive.M{
		"post_id": postID,
		"parent":  primitive.NilObjectID,
	}
	return col.CountDocuments(context.Background(), filter)
}

func ReplyCount(parentID primitive.ObjectID) (int64, error) {
	col := connection.GetDB().Comments
	filter := primitive.M{
		"parent": parentID,
	}
	return col.CountDocuments(context.Background(), filter)
}

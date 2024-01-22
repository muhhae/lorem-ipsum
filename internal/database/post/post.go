package post

import (
	"context"
	"time"

	"github.com/muhhae/lorem-ipsum/internal/database/connection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func FindOne(filter bson.M) (*Post, error) {
	var post Post
	err := connection.GetDB().Posts.FindOne(context.Background(), filter).Decode(&post)
	return &post, err
}

const (
	postLimit int64 = 2
)

func RetrievePosts(filter bson.M, iteration int64) ([]Post, error) {
	var posts []Post
	limits := postLimit
	skip := iteration * postLimit
	cursor, err := connection.GetDB().Posts.Find(context.Background(), filter, &options.FindOptions{
		Limit: &limits,
		Skip:  &skip,
		Sort:  bson.M{"createdAt": -1},
	})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &posts)
	return posts, err
}

func FindNewerOrOlder(anchor primitive.ObjectID, newer bool) ([]Post, error) {
	return nil, nil
}

func FindOlder(olderThan primitive.ObjectID) ([]Post, error) {
	col := connection.GetDB().Posts
	var posts []Post
	olderThanTime := time.Now()
	var olderThanPostId primitive.ObjectID

	if olderThan != primitive.NilObjectID {
		olderThanPost, err := FindOne(bson.M{"_id": olderThan})
		if err == nil {
			olderThanTime = olderThanPost.CreatedAt
			olderThanPostId = olderThanPost.ID
		}
	}
	filter := primitive.M{
		"$and": []interface{}{
			bson.M{
				"$or": []interface{}{
					bson.M{
						"createdAt": bson.M{"$lt": olderThanTime},
					},
					bson.M{
						"createdAt": olderThanTime,
						"_id": bson.M{
							"$lt": olderThanPostId,
						},
					},
				},
			},
			bson.M{
				"_id": bson.M{"$ne": olderThanPostId},
			},
		},
	}

	limit := postLimit
	cursor, err := col.Find(context.Background(), filter, &options.FindOptions{
		Limit: &limit,
		Sort:  bson.M{"createdAt": -1},
	})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &posts)
	return posts, err
}

func FindNewer(newerThan primitive.ObjectID) ([]Post, error) {
	col := connection.GetDB().Posts
	var posts []Post
	newerThanTime := time.Now()
	var newerThanPostId primitive.ObjectID

	if newerThan != primitive.NilObjectID {
		olderThanPost, err := FindOne(bson.M{"_id": newerThan})
		if err == nil {
			newerThanTime = olderThanPost.CreatedAt
			newerThanPostId = olderThanPost.ID
		}
	}
	filter := primitive.M{
		"$and": []interface{}{
			bson.M{
				"$or": []interface{}{
					bson.M{
						"createdAt": bson.M{"$gt": newerThanTime},
					},
					bson.M{
						"createdAt": newerThanTime,
						"_id": bson.M{
							"$gt": newerThanPostId,
						},
					},
				},
			},
			bson.M{
				"_id": bson.M{"$ne": newerThanPostId},
			},
		},
	}

	limit := postLimit
	cursor, err := col.Find(context.Background(), filter, &options.FindOptions{
		Limit: &limit,
		Sort:  bson.M{"createdAt": -1},
	})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &posts)
	return posts, err
}

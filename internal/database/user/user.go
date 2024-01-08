package user

import (
	"context"
	"encoding/json"

	"github.com/muhhae/lorem-ipsum/internal/database/connection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email    string             `json:"email" bson:"email,required"`
	Username string             `json:"username" bson:"username,required"`
	Password string             `json:"password" bson:"password,required"`
}

func FromJSON(u string) (*User, error) {
	var user *User
	err := json.Unmarshal([]byte(u), &user)
	if err != nil {
		return nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	return user, nil
}

func FindOne(filter User) (*User, error) {
	db := connection.GetDB()
	var user User
	res := db.Users.FindOne(context.Background(), filter)
	if res.Err() != nil {
		return nil, res.Err()
	}
	err := res.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FindAll(filter User) ([]*User, error) {
	db := connection.GetDB()
	var users []*User
	cur, err := db.Users.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (u *User) Save() error {
	db := connection.GetDB()
	if u.ID == primitive.NilObjectID {
		_, err := db.Users.InsertOne(context.Background(), u)
		return err
	}
	count, err := db.Users.CountDocuments(context.Background(), bson.M{"_id": u.ID})
	if err != nil {
		return err
	}
	if count > 0 {
		res := db.Users.FindOneAndReplace(context.Background(), bson.M{"_id": u.ID}, u)
		if res.Err() != nil {
			return res.Err()
		}
		return nil
	}
	_, err = db.Users.InsertOne(context.Background(), u)
	return err
}

func (u *User) ToJSON() (string, error) {
	j, err := json.Marshal(u)
	return string(j), err
}

func (u *User) Authenticate(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

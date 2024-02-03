package post

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	comment "github.com/muhhae/lorem-ipsum/internal/database/comment"
	imagemodel "github.com/muhhae/lorem-ipsum/internal/database/image"

	"github.com/muhhae/lorem-ipsum/internal/database/post"
	"github.com/muhhae/lorem-ipsum/internal/database/user"
	"github.com/muhhae/lorem-ipsum/internal/views/home"
	echotempl "github.com/muhhae/lorem-ipsum/pkg/echoTempl"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Upload(c echo.Context) error {
	authorID := c.Get("id").(primitive.ObjectID)
	if authorID == primitive.NilObjectID {
		return echo.ErrBadRequest
	}
	content := c.FormValue("content")
	if content == "" {
		return c.String(http.StatusBadRequest, "Post cannot be empty")
	}
	multipartForm, err := c.MultipartForm()
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid form")
	}
	images := multipartForm.File["images"]
	if len(images) == 0 {
		return c.String(http.StatusBadRequest, "no images uploaded")
	}
	if len(images) > 8 {
		return c.String(http.StatusBadRequest, "too many images")
	}
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/bmp":  true,
		"image/webp": true,
	}
	imageIDs := make([]primitive.ObjectID, len(images))
	counter := 0
	for img := range images {
		contentType := images[img].Header.Get("Content-Type")
		if !allowedTypes[contentType] {
			return c.String(http.StatusBadRequest, "invalid image type")
		}
		file, err := images[img].Open()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error processing image")
		}
		data, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error processing image")
		}
		img := imagemodel.Image{
			Owner: authorID,
			Data:  data,
		}
		imgID, err := img.Save()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error saving image")
		}
		imageIDs[counter] = imgID
		counter++
	}
	post := post.Post{
		AuthorID:  authorID,
		Content:   content,
		CreatedAt: time.Now(),
		ImageIDs:  imageIDs,
	}
	_, err = post.Save()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error saving post")
	}
	return c.NoContent(http.StatusOK)
}
func Delete(c echo.Context) error {
	ownerID := c.Get("id")
	if ownerID == "" {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}
	postIDstr := c.Param("id")
	if postIDstr == "" {
		return c.String(http.StatusBadRequest, "Invalid Post")
	}
	postID, err := primitive.ObjectIDFromHex(postIDstr)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	postData, err := post.FindOne(bson.M{
		"_id": postID,
	})
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	if postData.AuthorID != ownerID {
		return c.String(http.StatusUnauthorized, "Not yours man")
	}
	err = post.Delete(postID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func Default(c echo.Context) error {
	if c.QueryParam("olderThan") != "" {
		return Older(c)
	} else if c.QueryParam("newerThan") != "" {
		return Newer(c)
	}
	posts, err := post.FindOlder(primitive.NilObjectID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error retrieving posts")
	}
	postDatas, err := PostToPostdatas(posts, c.Get("id").(primitive.ObjectID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error retrieving posts")
	}
	if len(postDatas) == 0 {
		return echotempl.Templ(c, 200, home.EndOfFeed())
	}
	return echotempl.Templ(c, 200, home.ManyPost(postDatas, home.ManyPostTypeBoth))
}

func Older(c echo.Context) error {
	olderThan := c.QueryParam("olderThan")
	olderThanID, err := primitive.ObjectIDFromHex(olderThan)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid olderThan id")
	}
	posts, err := post.FindOlder(olderThanID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error retrieving posts")
	}
	postDatas, err := PostToPostdatas(posts, c.Get("id").(primitive.ObjectID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error retrieving posts")
	}
	if len(postDatas) == 0 {
		return echotempl.Templ(c, 200, home.EndOfFeed())
	}
	return echotempl.Templ(c, 200, home.ManyPost(postDatas, home.ManyPostTypeOlder))

}

func Newer(c echo.Context) error {
	newerThan := c.QueryParam("newerThan")
	newerThanID, err := primitive.ObjectIDFromHex(newerThan)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid olderThan id")
	}
	posts, err := post.FindNewer(newerThanID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error retrieving posts")
	}
	postDatas, err := PostToPostdatas(posts, c.Get("id").(primitive.ObjectID))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error retrieving posts")
	}
	if len(postDatas) == 0 {
		return c.String(http.StatusNoContent, "No newer posts")
	}
	return echotempl.Templ(c, 200, home.ManyPost(postDatas, home.ManyPostTypeNewer))
}

func PostToPostdatas(posts []post.Post, userID primitive.ObjectID) ([]home.PostData, error) {
	postDatas := make([]home.PostData, len(posts))
	for i, p := range posts {
		owner, err := user.FindOne(bson.M{"_id": p.AuthorID})
		if err != nil {
			return nil, errors.New("error retrieving posts")
		}
		images := make([]string, len(p.ImageIDs))
		for i, imgID := range p.ImageIDs {
			images[i] = fmt.Sprintf("/api/v1/image/%s", imgID.Hex())
		}
		likeCount, err := post.CountReaction(p.ID)
		if err != nil {
			return nil, errors.New("error retrieving posts")
		}
		userID := userID
		var v int
		if userID == primitive.NilObjectID {
			v = 0
		} else {
			v, err = post.GetReaction(p.ID, userID)
			if err != nil {
				return nil, errors.New("error retrieving posts")
			}
		}
		commentCount, err := comment.CommentCount(p.ID)
		if err != nil {
			commentCount = 0
		}
		postData := home.PostData{
			PostID:       p.ID.Hex(),
			Username:     owner.Username,
			Content:      p.Content,
			ImgSrc:       images,
			CommentCount: int(commentCount),
			ReactStruct: home.ReactData{
				PostID:    p.ID.Hex(),
				LikeCount: int(likeCount),
				Value:     v,
			},
		}
		postDatas[i] = postData
	}
	return postDatas, nil
}

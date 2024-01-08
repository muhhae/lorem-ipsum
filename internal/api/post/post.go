package post

import (
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	imagemodel "github.com/muhhae/lorem-ipsum/internal/database/image"
	"github.com/muhhae/lorem-ipsum/internal/database/post"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Upload(c echo.Context) error {
	authorID := c.Get("id").(primitive.ObjectID)
	if authorID == primitive.NilObjectID {
		return echo.ErrBadRequest
	}

	content := c.FormValue("content")
	imageData, err := c.FormFile("image")
	if err != nil {
		return echo.ErrBadRequest
	}

	file, err := imageData.Open()
	if err != nil {
		return echo.ErrBadRequest
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return echo.ErrBadRequest
	}

	img := imagemodel.Image{
		Owner: authorID,
		Data:  data,
	}

	imgID, err := img.Save()
	if err != nil {
		return echo.ErrBadRequest
	}

	post := post.Post{
		AuthorID:  authorID,
		Content:   content,
		CreatedAt: time.Now(),
		ImageIDs:  []primitive.ObjectID{imgID},
	}

	_, err = post.Save()
	if err != nil {
		return echo.ErrBadRequest
	}
	c.JSON(http.StatusOK, map[string]string{"message": "success"})
	return nil
}

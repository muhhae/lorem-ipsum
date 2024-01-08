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

func Upload(c echo.Context) {
	authorID := c.Get("id").(primitive.ObjectID)
	if authorID == primitive.NilObjectID {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid author"})
		return
	}

	content := c.FormValue("content")
	imageData, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	file, err := imageData.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	img := imagemodel.Image{
		Owner: authorID,
		Data:  data,
	}

	imgID, err := img.Save()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	post := post.Post{
		AuthorID:  authorID,
		Content:   content,
		CreatedAt: time.Now(),
		ImageIDs:  []primitive.ObjectID{imgID},
	}

	_, err = post.Save()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]string{"message": "success"})
}

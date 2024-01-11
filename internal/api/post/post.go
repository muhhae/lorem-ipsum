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

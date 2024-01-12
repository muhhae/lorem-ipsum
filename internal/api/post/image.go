package post

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	imagemodel "github.com/muhhae/lorem-ipsum/internal/database/image"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetImage(c echo.Context) error {
	fmt.Println("GetImage", c.Param("id"))
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid image ID")
	}
	img, err := imagemodel.FindOne(bson.M{"_id": id})
	if err != nil || img == nil {
		fmt.Println("1", err)
		return c.String(404, "Image not found")
	}
	if img.Data == nil {
		fmt.Println("2")
		return c.String(404, "Image not found")
	}

	contentType := http.DetectContentType(img.Data)
	if contentType == "" {
		fmt.Println("3")
		return c.String(404, "Image not found")
	}
	// c.Response().Header().Set("Cache-Control", "max-age=31536000")
	// c.Response().Header().Set("Expires", time.Now().AddDate(1, 0, 0).Format(http.TimeFormat))
	// c.Response().Header().Set("Last-Modified", time.Now().Format(http.TimeFormat))
	// c.Response().Header().Set("ETag", fmt.Sprintf("%x", md5.Sum(img.Data)))
	return c.Blob(200, contentType, img.Data)
}

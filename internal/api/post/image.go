package post

import (
	"net/http"

	"github.com/labstack/echo/v4"
	imagemodel "github.com/muhhae/lorem-ipsum/internal/database/image"
	"go.mongodb.org/mongo-driver/bson"
)

func GetImage(c echo.Context) error {
	img, err := imagemodel.FindOne(bson.M{"_id": c.Param("id")})
	if err != nil || img == nil {
		return c.String(404, "Image not found")
	}
	if img.Data == nil {
		return c.String(404, "Image not found")
	}

	contentType := http.DetectContentType(img.Data)
	if contentType == "" {
		return c.String(404, "Image not found")
	}
	// c.Response().Header().Set("Cache-Control", "max-age=31536000")
	// c.Response().Header().Set("Expires", time.Now().AddDate(1, 0, 0).Format(http.TimeFormat))
	// c.Response().Header().Set("Last-Modified", time.Now().Format(http.TimeFormat))
	// c.Response().Header().Set("ETag", fmt.Sprintf("%x", md5.Sum(img.Data)))
	return c.Blob(200, contentType, img.Data)
}

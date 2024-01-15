package comment

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/muhhae/lorem-ipsum/internal/database/comment"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SendComment(c echo.Context) error {
	id := c.Get("id").(primitive.ObjectID)
	if id == primitive.NilObjectID {
		return c.String(http.StatusUnauthorized, "You must be logged in to comment")
	}
	postID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid post id")
	}
	content := c.FormValue("content")
	if content == "" {
		return c.String(http.StatusBadRequest, "Comment cannot be empty")
	}
	parentID, err := primitive.ObjectIDFromHex(c.FormValue("parent"))
	if err != nil {
		parentID = primitive.NilObjectID
	}
	comment := comment.Comment{
		PostID:    postID,
		UserID:    id,
		Content:   content,
		Parent:    parentID,
		CreatedAt: time.Now(),
	}
	_, err = comment.Save()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error saving comment")
	}
	return c.String(http.StatusOK, "Comment saved")
}

func GetPostComment(c echo.Context) error {
	postID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid post id")
	}
	iteration, err := strconv.Atoi(c.QueryParam("iteration"))
	if err != nil {
		iteration = 0
	}
	comments, err := comment.RetrieveDefault(postID, iteration)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error getting comments")
	}
	return c.JSON(http.StatusOK, comments)
}

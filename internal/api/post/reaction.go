package post

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/muhhae/lorem-ipsum/internal/database/post"
	"github.com/muhhae/lorem-ipsum/internal/views/home"
	"github.com/muhhae/lorem-ipsum/internal/views/util"
	echotempl "github.com/muhhae/lorem-ipsum/pkg/echoTempl"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func React(c echo.Context) error {
	postID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil || postID == primitive.NilObjectID {
		return c.String(http.StatusBadRequest, "Invalid post ID")
	}
	userID := c.Get("id").(primitive.ObjectID)
	if userID == primitive.NilObjectID {
		return c.String(http.StatusBadRequest, "You need to be logged in to react")
	}
	value, err := strconv.Atoi(c.QueryParam("value"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid reaction value")
	}
	if value < -1 || value > 1 {
		return c.String(http.StatusBadRequest, "Invalid reaction value")
	}
	react := post.Reaction{
		PostID: postID,
		UserID: userID,
		Value:  value,
	}
	_, err = react.Save()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	return c.String(http.StatusOK, "Reaction saved")
}

func ReactionCount(c echo.Context) error {
	postID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil || postID == primitive.NilObjectID {
		return c.String(http.StatusBadRequest, "Invalid post ID")
	}
	count, err := post.CountReaction(postID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	return c.String(http.StatusOK, util.Format(int(count)))
}

func MyReaction(c echo.Context) error {
	postID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil || postID == primitive.NilObjectID {
		return c.String(http.StatusBadRequest, "Invalid post ID")
	}
	userID := c.Get("id").(primitive.ObjectID)
	if userID == primitive.NilObjectID {
		return c.String(http.StatusBadRequest, "You need to be logged in to react")
	}
	reaction, err := post.GetReaction(postID, userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	return c.String(http.StatusOK, strconv.Itoa(reaction))
}

func GetReaction(c echo.Context) error {
	postID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil || postID == primitive.NilObjectID {
		return c.String(http.StatusBadRequest, "Invalid post ID")
	}
	var reaction int
	userID := c.Get("id").(primitive.ObjectID)
	if userID == primitive.NilObjectID {
		reaction = 0
	} else {
		reaction, err = post.GetReaction(postID, userID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal server error")
		}
	}
	reactCount, err := post.CountReaction(postID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	data := home.ReactData{
		PostID:    postID.Hex(),
		LikeCount: int(reactCount),
		Value:     reaction,
	}
	return echotempl.Templ(c, 200, home.ReactSection(data))
}

package comment

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/muhhae/lorem-ipsum/internal/database/comment"
	"github.com/muhhae/lorem-ipsum/internal/database/user"
	"github.com/muhhae/lorem-ipsum/internal/views/home"
	"github.com/muhhae/lorem-ipsum/internal/views/util"
	echotempl "github.com/muhhae/lorem-ipsum/pkg/echoTempl"
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
	parentID, err := primitive.ObjectIDFromHex(c.QueryParam("replying"))
	if err != nil || parentID == primitive.NilObjectID {
		parentID = primitive.NilObjectID
	}
	content := c.FormValue("content")
	if content == "" {
		return c.String(http.StatusBadRequest, "Comment cannot be empty")
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
	var after primitive.ObjectID
	afterP := c.QueryParam("after")
	if afterP == "" {
		after = primitive.NilObjectID
	} else {
		after, err = primitive.ObjectIDFromHex(afterP)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid after id")
		}
	}
	comments, err := comment.RetrieveAll(postID, primitive.NilObjectID, after)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error getting comments")
	}
	var commentDatas []home.CommentData
	for _, cm := range comments {
		user, err := user.FindOne(primitive.M{"_id": cm.UserID})
		if err != nil || user.ID == primitive.NilObjectID {
			continue
		}
		replyCount, err := comment.ReplyCount(cm.ID)
		if err != nil {
			replyCount = 0
		}
		commentData := home.CommentData{
			PostID:     postID.Hex(),
			CommentID:  cm.ID.Hex(),
			Content:    cm.Content,
			Username:   user.Username,
			ReplyCount: int(replyCount),
			Time:       cm.CreatedAt.Format(time.RFC3339),
		}
		commentDatas = append(commentDatas, commentData)
	}
	if len(commentDatas) == 0 {
		return c.String(http.StatusNotFound, "No comments found")
	}

	return echotempl.Templ(c, 200, home.LoadedComment(commentDatas, true))

}

func GetReply(c echo.Context) error {
	parentID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid parent id")
	}
	comments, err := comment.RetrieveAll(primitive.NilObjectID, parentID, primitive.NilObjectID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error getting comments")
	}
	var commentDatas []home.CommentData
	for _, cm := range comments {
		user, err := user.FindOne(primitive.M{"_id": cm.UserID})
		if err != nil || user.ID == primitive.NilObjectID {
			continue
		}
		replyCount, err := comment.ReplyCount(cm.ID)
		if err != nil {
			replyCount = 0
		}
		commentData := home.CommentData{
			PostID:     cm.PostID.Hex(),
			CommentID:  cm.ID.Hex(),
			Content:    cm.Content,
			Username:   user.Username,
			ReplyCount: int(replyCount),
		}
		commentDatas = append(commentDatas, commentData)
	}

	return echotempl.Templ(c, 200, home.LoadedComment(commentDatas, false))
}

func GetCommentCount(c echo.Context) error {
	postID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid post id")
	}
	count, err := comment.CommentCount(postID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error getting comment count")
	}
	return c.String(200, util.Format(int(count)))
}

func GetReplyCount(c echo.Context) error {
	parentID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid parent id")
	}
	count, err := comment.ReplyCount(parentID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error getting reply count")
	}
	return c.String(200, util.Format(int(count)))
}

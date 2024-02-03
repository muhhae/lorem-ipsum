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

func GetComment(c echo.Context) error {
	postID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid post id")
	}
	after := primitive.NilObjectID
	afterP := c.QueryParam("after")
	if afterP != "" {
		after, err = primitive.ObjectIDFromHex(afterP)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid after id : "+afterP)
		}
	}
	parent := primitive.NilObjectID
	parentP := c.QueryParam("parent")
	if parentP != "" {
		parent, err = primitive.ObjectIDFromHex(parentP)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid parent id : "+parentP)
		}
	}
	comments, err := comment.RetrieveAll(postID, parent, after)
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
			ParentID:   cm.Parent.Hex(),
			CommentID:  cm.ID.Hex(),
			Content:    cm.Content,
			Username:   user.Username,
			ReplyCount: int(replyCount),
		}
		commentDatas = append(commentDatas, commentData)
	}
	if len(commentDatas) == 0 {
		return c.String(http.StatusNotFound, "No comments found")
	}
	return echotempl.Templ(c, 200, home.LoadedComment(commentDatas))
}

func GetCommentCount(c echo.Context) error {
	postID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid post id")
	}
	parentID, err := primitive.ObjectIDFromHex(c.QueryParam("parent"))
	if err != nil {
		parentID = primitive.NilObjectID
	}
	if parentID == primitive.NilObjectID {
		count, err := comment.CommentCount(postID)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error getting comment count")
		}
		return c.String(200, util.Format(int(count)))
	}
	count, err := comment.ReplyCount(parentID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error getting reply count")
	}
	return c.String(200, util.Format(int(count)))
}

package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateCommentRequest struct {
	Content string `json:"content"`
	PostId  string `json:"postId"`
}
type DeleteCommentRequest struct {
	CommentId string `json:"commentId"`
}
type UpdateCommentRequest struct {
	CommentId string `json:"commentId"`
	Content   string `json:"content"`
}
type GetPostsCommentsRequest struct {
	PostId string `json:"postId"`
}

func (c *AppController) CreatePostComment(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(string)

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"message": "CreateComment" + userId,
	})

}
func (c *AppController) DeletePostComment(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(string)

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"message": "DeleteComment" + userId,
	})

}
func (c *AppController) UpdatePostComment(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(string)

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"message": "UpdateComment" + userId,
	})

}
func (c *AppController) GetPostsComments(ctx *gin.Context) {
	userId := ctx.MustGet("userId").(string)
	postId := ctx.Param("id")

	ctx.IndentedJSON(http.StatusOK, gin.H{
		"message": "GetPostsComments" + userId + " " + postId,
	})

}

package main

import (
	"net/http"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllPostsRest(c *gin.Context) {
	posts, err := GetAllPostsWithoutComments()
	if err != nil {
		fmt.Println(err)
		return
	}
	c.IndentedJSON(http.StatusOK, posts)
}

func GetPostByIDRest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Invalid ID %v: Not an integer", idStr)})
		return
	}
	post, err := GetPostByID(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	comments, err := GetCommentsByPostID(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	post.Comments = comments
	c.IndentedJSON(http.StatusOK, post)
}

func PostNewPost(c *gin.Context) {
	var post Post
	if err := c.BindJSON(&post); err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON for post"})
		return
	}
	_, err := AddPost(&post)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.IndentedJSON(http.StatusCreated, post)
}

// TODO: Remove repetition
func PostNewComment(c *gin.Context) {
	var comment Comment
	if err := c.BindJSON(&comment); err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON for comment"})
		return
	}
	_, err := AddComment(&comment)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.IndentedJSON(http.StatusCreated, comment)
}
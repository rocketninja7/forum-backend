package main

import (
	"net/http"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllPostsRest(c *gin.Context) {
	posts, err := GetAllPostsWithUser()
	if err != nil {
		fmt.Println(err)
		return
	}
	posts, err = GetAllTagsForPosts(posts)
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
	tags, err := GetTagsByPostID(id)
	if err != nil {
		fmt.Println(err)
		return
	}
	post.Tags = tags
	c.IndentedJSON(http.StatusOK, post)
}

func PostNewPostRest(c *gin.Context) {
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
func PostNewCommentRest(c *gin.Context) {
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

func GetPostsByTagRest(c *gin.Context) {
	/*
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Invalid ID %v: Not an integer", idStr)})
		return
	}
	*/
}
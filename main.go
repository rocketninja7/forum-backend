package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	_ "github.com/lib/pq"
)

type User struct {
	Id int64 `json:"id"`
	Username string `json:"username"`
}

type Post struct {
	Id int64 `json:"id"`
	Poster User `json:"poster"`
	TimeCreated time.Time `json:"time_created"`
	Title string `json:"title"`
	Content string `json:"content"`
	Comments []Comment `json:"comments"`
	Tags []Tag `json:"tags"`
}

type Comment struct {
	Id int64 `json:"id"`
	TimeCreated time.Time `json:"time_created"`
	PostId int64 `json:"post"`
	Commenter User `json:"commenter"`
	Content string `json:"content"`
}

type Tag struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
}

var db *sql.DB

func main() {
	connStr := "dbname=forum user=postgres password=1234 sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("Opened!")

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{"Content-type"},
		AllowMethods: []string{"DELETE", "PUT"},
	}))
	
	router.GET("/", GetAllPostsRest)
	router.GET("/post/:id/", GetPostByIDRest)
	// router.POST("/login/")
	// router.GET("/user/:id/")
	// router.GET("/tag/:id/")
	router.POST("/newpost/", PostNewPostRest)
	router.POST("/newcomment/", PostNewCommentRest)
	router.DELETE("/post/:id/", DeletePostByIdRest)
	router.DELETE("/comment/:id/", DeleteCommentByIdRest)
	router.PUT("post/:id/", PutPostByIdRest)
	router.PUT("comment/:id/", PutCommentByIdRest)

	router.Run("localhost:8080")
}





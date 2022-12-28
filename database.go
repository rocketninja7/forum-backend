package main

import (
	"database/sql"
	"fmt"
	"time"
)

func GetAllPostsWithoutComments() ([]Post, error) {
	var posts []Post

	rows, err := db.Query(
		"SELECT post.\"id\", post.poster, post.time_created, post.title, post.\"content\", \"user\".username " + 
		"FROM post JOIN \"user\" ON post.poster=\"user\".\"id\"")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		var user User
		if err := rows.Scan(&post.Id, &user.Id, &post.TimeCreated, &post.Title, &post.Content, &user.Username); err != nil {
			return nil, err
		}
		post.Poster = user
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetPostByID(id int64) (Post, error) {
	var post Post
	var user User

	row := db.QueryRow(
		"SELECT post.\"id\", post.poster, post.time_created, post.title, post.\"content\", \"user\".username " +
		"FROM post JOIN \"user\" ON post.poster=\"user\".\"id\" " + 
		"WHERE post.\"id\"=$1", id)
	if err := row.Scan(&post.Id, &user.Id, &post.TimeCreated, &post.Title, &post.Content, &user.Username); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Post not found!")
		}
		return post, err
	}
	post.Poster = user
	return post, nil
}

func GetCommentsByPostID(id int64) ([]Comment, error) {
	var comments []Comment

	rows, err := db.Query(
		"SELECT \"comment\".\"id\", \"comment\".time_created, \"comment\".post, \"comment\".commenter, \"comment\".\"content\", \"user\".username " + 
		"FROM \"comment\" JOIN \"user\" ON \"comment\".commenter=\"user\".\"id\" " + 
		"WHERE \"comment\".post=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		var user User
		if err := rows.Scan(&comment.Id, &comment.TimeCreated, &comment.Post, &user.Id, &comment.Content, &user.Username); err != nil {
			return nil, err
		}
		comment.Commenter = user
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func AddPost(post *Post) (int64, error) {
	err := db.QueryRow(
		"INSERT INTO post (poster, time_created, title, \"content\") VALUES ($1, $2, $3, $4) RETURNING \"id\"", post.Poster.Id, time.Now(), post.Title, post.Content).Scan(&post.Id)
	if err != nil {
		return 0, err
	}
	return post.Id, nil
}
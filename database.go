package main

import (
	"database/sql"
	"fmt"
	"time"
)

func GetAllPostsWithUser() ([]Post, error) {
	var posts []Post

	rows, err := db.Query(
		"SELECT post.\"id\", post.poster, post.time_created, post.title, post.\"content\", \"user\".username " + 
		"FROM post JOIN \"user\" ON post.poster=\"user\".\"id\" ORDER BY post.time_created")
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

func GetAllTagsForPosts(posts []Post) ([]Post, error) {
	postsById := make(map[int64]*Post)
	for idx := range posts {
		posts[idx].Tags = make([]Tag, 0)
		postsById[posts[idx].Id] = &posts[idx]
	}

	rows, err := db.Query(
		"SELECT post_tag.post_id, post_tag.tag_id, tag.\"name\" " + 
		"FROM post_tag JOIN tag ON post_tag.tag_id=tag.\"id\" ORDER BY post_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var postId int64
		var tag Tag
		if err := rows.Scan(&postId, &tag.Id, &tag.Name); err != nil {
			return nil, err
		}
		currPost, ok := postsById[postId]
		if ok {
			currPost.Tags = append(currPost.Tags, tag)
		}
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
	comments := make([]Comment, 0)

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
		if err := rows.Scan(&comment.Id, &comment.TimeCreated, &comment.PostId, &user.Id, &comment.Content, &user.Username); err != nil {
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

func GetTagsByPostID(id int64) ([]Tag, error) {
	tags := make([]Tag, 0)

	rows, err := db.Query(
		"SELECT tag.\"id\", tag.\"name\" " + 
		"FROM post_tag JOIN tag ON post_tag.tag_id=tag.\"id\" " + 
		"WHERE post_tag.post_id=$1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.Id, &tag.Name); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func AddPost(post *Post) (int64, error) {
	err := db.QueryRow(
		"INSERT INTO post (poster, time_created, title, \"content\") VALUES ($1, $2, $3, $4) RETURNING \"id\"", post.Poster.Id, time.Now(), post.Title, post.Content).Scan(&post.Id)
	if err != nil {
		return 0, err
	}
	return post.Id, nil
}

func AddComment(comment *Comment) (int64, error) {
	err := db.QueryRow(
		"INSERT INTO comment (post, commenter, time_created, \"content\") VALUES ($1, $2, $3, $4) RETURNING \"id\"", comment.PostId, comment.Commenter.Id, time.Now(), comment.Content).Scan(&comment.Id)
	if err != nil {
		return 0, err
	}
	return comment.Id, nil
}

func DeletePost(id int64) (bool, error) {
	res, err := db.Exec("DELETE FROM post WHERE \"id\"=$1", id)
	if err != nil {
		return false, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func DeleteComment(id int64) (bool, error) {
	res, err := db.Exec("DELETE FROM \"comment\" WHERE \"id\"=$1", id)
	if err != nil {
		return false, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func UpdatePost(post Post) (bool, error) {
	res, err := db.Exec("UPDATE post SET title=$1, content=$2 WHERE \"id\"=$3", post.Title, post.Content, post.Id)
	if err != nil {
		return false, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func UpdateComment(comment Comment) (bool, error) {
	res, err := db.Exec("UPDATE \"comment\" SET content=$1 WHERE \"id\"=$2", comment.Content, comment.Id)
	if err != nil {
		return false, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

// TODO: Remove repetition, currently not used
func GetPostsByTagID(id int64) ([]Post, error) {
	var posts []Post

	rows, err := db.Query(
		"SELECT post.\"id\", post.poster, post.time_created, post.title, post.\"content\", \"user\".username " + 
		"FROM post_tag JOIN post ON post_tag.post_id=post.\"id\" JOIN \"user\" ON post.poster=\"user\".\"id\" " + 
		"WHERE post_tag.tag_id=$1", id)
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
# Forum Backend
A simple backend built with REST API endpoints.
To be used with https://github.com/rocketninja7/forum-frontend.

## Endpoints:
### GET
- "/": Get all posts in the forum.
- "/post/:id/": Get the post with the specified id.
### POST
- "/newpost/": Create a new post.
- "/newcomment/": Create a new comment.
### DELETE
- "/post/:id/": Delete the post with the specified id.
- "/comment/:id/": Delete the comment with the specified id.
### PUT
- "/post/:id/": Update a post with the specified id.
- "/comment/:id/": Update a comment with the specified id.
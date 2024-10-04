package postgres

import (
	"fmt"

	"github.com/alaalser/goreddit"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// func NewPostStore(db *sqlx.DB) *PostStore {
// 	return &PostStore{
// 		DB: db,
// 	}
// }

type PostStore struct {
	*sqlx.DB
}

func (s *PostStore) Post(id uuid.UUID) (goreddit.Post, error) {
	var post goreddit.Post
	if err := s.Get(&post, `SELECT * FROM posts WHERE id = %1`, id); err != nil {
		return goreddit.Post{}, fmt.Errorf("error getting post %w", err)
	}
	return post, nil
}

func (s *PostStore) PostsByThread(threadID uuid.UUID) ([]goreddit.Post, error) {
	var posts []goreddit.Post
	if err := s.Select(&posts, `SELECT * FROM posts WHERE thread_id=$1`, threadID); err != nil {
		return []goreddit.Post{}, fmt.Errorf("error getting posts %w", err)
	}
	return posts, nil
}

func (s *PostStore) CreatePost(post *goreddit.Post) error {
	if err := s.Get(post, `INSERT INTO posts VALUES($1,$2,$3,$4,$5) RETURNING *`, post.ID, post.ThreadID, post.Title, post.Content, post.Votes); err != nil {
		return fmt.Errorf("error creating pots %w", err)
	}
	return nil
}

func (s *PostStore) UpdatePost(post *goreddit.Post) error {
	if err := s.Get(post, `UPDATE posts SET thread_id = $1, title=$2, content=$3,votes=$4 WHERE id=$5) RETURNING *`, post.ThreadID, post.Title, post.Content, post.Votes, post.ID); err != nil {
		return fmt.Errorf("error creating pots %w", err)
	}
	return nil
}

func (s *PostStore) DeletePost(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM posts where id=$1 RETURNING *`, id); err != nil {
		return fmt.Errorf("error deleting posts")
	}
	return nil
}

package postgres

import (
	"fmt"

	"github.com/alaalser/goreddit"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// func NewCommentStore(db *sqlx.DB) *CommentStore {
// 	return &CommentStore{
// 		DB: db,
// 	}
// }

type CommentStore struct {
	*sqlx.DB
}

func (s *CommentStore) Comment(id uuid.UUID) (goreddit.Comment, error) {
	var comment goreddit.Comment
	if err := s.Get(&comment, `SELECT * FROM comments WHERE id=$1`, id); err != nil {
		return goreddit.Comment{}, fmt.Errorf("error getting comment")
	}
	return comment, nil
}

func (s *CommentStore) CommentsByPost(postID uuid.UUID) ([]goreddit.Comment, error) {
	var comments []goreddit.Comment
	if err := s.Select(&comments, `SELECT * FROM comments WHERE post_id=$1`, postID); err != nil {
		return []goreddit.Comment{}, fmt.Errorf("error getting comments")
	}
	return comments, nil
}

func (s *CommentStore) CreateComment(comment *goreddit.Comment) error {
	if err := s.Get(comment, `INSERT INTO comments VALUES($1,$2,$3,$4) RETURNING *`, comment.ID, comment.PostID, comment.Content, comment.Votes); err != nil {
		return fmt.Errorf("error creating comment %w", err)
	}
	return nil
}

func (s *CommentStore) UpdateComment(comment *goreddit.Comment) error {
	if err := s.Get(comment, `UPDATE comments SET post_id = $1, content$2, votes = $3 where id=$4) RETURNING *`, comment.PostID, comment.Content, comment.Votes, comment.ID); err != nil {
		return fmt.Errorf("error creating comment %w", err)
	}
	return nil
}

func (s *CommentStore) DeleteComment(id uuid.UUID) error {
	if _, err := s.Exec("DELETE FROM comments WHERE id= $1", id); err != nil {
		return fmt.Errorf("error deleting comment %w", err)
	}
	return nil
}

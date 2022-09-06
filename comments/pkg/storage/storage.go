package storage

import (
	"gopkg.in/guregu/null.v4/zero"
)

type Comment struct {
	ID              int      `json:"ID,omitempty"`
	NewsID          int      `json:"newsID,omitempty"`
	ParentCommentID zero.Int `json:"parentCommentID,omitempty"`
	Content         string   `json:"content,omitempty"`
	PubTime         int64    `json:"pubTime,omitempty"`
}

type Interface interface {
	AllComments(newsID int) ([]Comment, error)
	AddComment(Comment) error
	DeleteComment(Comment) error
}

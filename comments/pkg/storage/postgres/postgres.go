package postgres

import (
	"comments/pkg/storage"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type postgres struct {
	db *pgxpool.Pool
}

func New(url string) (*postgres, error) {
	log.Println("Try connect to db_comments")
	for {
		_, err := pgxpool.Connect(context.Background(), url)
		if err == nil {
			break
		}
	}
	db, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}
	log.Println("The connection was successful")
	return &postgres{db: db}, nil
}

func (p *postgres) AllComments(newsID int) ([]storage.Comment, error) {
	rows, err := p.db.Query(context.Background(), "SELECT * FROM comments WHERE news_id = $1;", newsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []storage.Comment
	for rows.Next() {
		var c storage.Comment
		err = rows.Scan(&c.ID, &c.NewsID, &c.ParentCommentID, &c.Content, &c.PubTime)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (p *postgres) AddComment(c storage.Comment) error {
	if c.NewsID == 0 {
		var newsId int
		row := p.db.QueryRow(context.Background(), "SELECT news_id FROM comments WHERE id =$1;", c.ParentCommentID)
		err := row.Scan(&newsId)
		if err != nil {
			return err
		}
		_, err = p.db.Exec(context.Background(),
			"INSERT INTO comments (news_id,parent_comment_id,content) VALUES ($1,$2,$3);", newsId, c.ParentCommentID, c.Content)
		if err != nil {
			return err
		}
		return nil
	} else {
		_, err := p.db.Exec(context.Background(),
			"INSERT INTO comments (news_id,content) VALUES ($1,$2);", c.NewsID, c.Content)
		if err != nil {
			return err
		}
		return nil
	}

}

func (p *postgres) DeleteComment(c storage.Comment) error {
	_, err := p.db.Exec(context.Background(),
		"DELETE FROM comments WHERE id=$1;", c.ID)
	if err != nil {
		return err
	}
	return nil
}

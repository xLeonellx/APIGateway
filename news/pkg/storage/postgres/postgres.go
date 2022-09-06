package postgres

import (
	"GoNews/pkg/storage"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type postgres struct {
	db *pgxpool.Pool
}

// Создание объекта DB
func New(url string) (*postgres, error) {
	log.Println("Try connect to db_news")
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

func (p *postgres) PostsMany(posts []storage.Post) error {
	for _, post := range posts {
		err := p.AddPost(post)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *postgres) PostSearchILIKE(pattern string, limit, offset int) ([]storage.Post, storage.Pagination, error) {
	pattern = "%" + pattern + "%"

	pagination := storage.Pagination{
		Page:  offset/limit + 1,
		Limit: limit,
	}
	row := p.db.QueryRow(context.Background(), "SELECT count(*) FROM posts WHERE title ILIKE $1;", pattern)
	err := row.Scan(&pagination.NumOfPages)

	if pagination.NumOfPages%limit > 0 {
		pagination.NumOfPages = pagination.NumOfPages/limit + 1
	} else {
		pagination.NumOfPages /= limit
	}

	if err != nil {
		return nil, storage.Pagination{}, err
	}

	rows, err := p.db.Query(context.Background(), "SELECT * FROM posts WHERE title ILIKE $1 ORDER BY pubtime DESC LIMIT $2 OFFSET $3;", pattern, limit, offset)
	if err != nil {
		return nil, storage.Pagination{}, err
	}
	defer rows.Close()
	var posts []storage.Post
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.PubTime, &p.Link)
		if err != nil {
			return nil, storage.Pagination{}, err
		}
		posts = append(posts, p)
	}
	return posts, pagination, rows.Err()
}

func (p *postgres) PostDetal(id int) (storage.Post, error) {
	row := p.db.QueryRow(context.Background(), "SELECT * FROM posts WHERE id =$1;", id)
	var post storage.Post
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.PubTime, &post.Link)
	if err != nil {
		return storage.Post{}, err
	}
	return post, nil
}

func (p *postgres) Posts(limit, offset int) ([]storage.Post, error) {
	rows, err := p.db.Query(context.Background(), "SELECT * FROM posts ORDER BY pubtime DESC LIMIT $1 OFFSET $2;", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []storage.Post
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.PubTime, &p.Link)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (p *postgres) AddPost(post storage.Post) error {
	_, err := p.db.Exec(context.Background(),
		"INSERT INTO posts (title, content, pubtime, link) VALUES ($1,$2, $3, $4);", post.Title, post.Content, post.PubTime, post.Link)
	if err != nil {
		return err
	}
	return nil
}

func (p *postgres) UpdatePost(post storage.Post) error {
	_, err := p.db.Exec(context.Background(),
		"UPDATE posts "+
			"SET title = $1, "+
			"content = $2, "+
			"pubtime = $3,"+
			"link = $4 "+
			"WHERE id = $5", post.Title, post.Content, post.PubTime, post.Link, post.ID)
	if err != nil {
		return err
	}
	return nil
}

func (p *postgres) DeletePost(post storage.Post) error {
	_, err := p.db.Exec(context.Background(),
		"DELETE FROM posts WHERE id=$1;", post.ID)
	if err != nil {
		return err
	}
	return nil
}

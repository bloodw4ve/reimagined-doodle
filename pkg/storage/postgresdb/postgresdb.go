package postgresdb

import (
	"context"
	"mod36a41/pkg/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

// connect to postgresdb
func New(conn string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), conn)

	if err != nil {
		return nil, err
	}

	s := Store{
		db: db,
	}

	return &s, nil
}

// closes connection to postgresdb
func (s *Store) Close() {
	s.db.Close()
}

// creates an array
func (s *Store) StoreNews(news []storage.Post) error {
	for _, p := range news {
		_, err := s.db.Exec(context.Background(), `
		INSERT INTO news(title, content, pub_time, link)
		VALUES ($1, $2, $3, $4)`,
			p.Title,
			p.Content,
			p.PubTime,
			p.Link)
		if err != nil {
			return err
		}
	}
	return nil
}

// returns n number of news posts based on date from db
func (s *Store) LastNews(n int) ([]storage.Post, error) {
	if n == 0 {
		n = 10
	}

	rows, err := s.db.Query(context.Background(), `
		SELECT id, title, content, pub_time, link FROM news
		ORDER BY pub_time DESC
		LIMIT $1
		`, n)

	if err != nil {
		return nil, err
	}

	var news []storage.Post
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}
		news = append(news, p)
	}
	return news, rows.Err()
}

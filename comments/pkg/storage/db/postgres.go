package db

import (
	"comments/pkg/storage"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CommentDB struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, dbUrl string) (*CommentDB, error) {
	for {
		_, err := pgxpool.Connect(ctx, dbUrl)
		if err == nil {
			break
		}
	}
	db, err := pgxpool.Connect(ctx, dbUrl)
	if err != nil {
		return nil, err
	}
	c := CommentDB{
		db: db,
	}
	return &c, nil
}

func (c *CommentDB) GetComments(IdNews int) ([]storage.Comment, error) {
	rows, err := c.db.Query(context.Background(), `SELECT * FROM comments WHERE id_news = $1`, IdNews)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comms []storage.Comment
	for rows.Next() {
		var com storage.Comment
		err = rows.Scan(&com.ID, &com.IdNews, &com.Content, &com.PubTime)
		if err != nil {
			return nil, err
		}
		comms = append(comms, com)
	}
	return comms, rows.Err()
}

func (c *CommentDB) AddComment(s storage.Comment) error {
	err := c.db.QueryRow(context.Background(),
		`INSERT INTO comments (id_news, content) VALUES 
             ($1, $2);`, s.IdNews, s.Content).Scan()
	if err != nil {
		return err
	}
	return nil
}

func (c *CommentDB) RemoveComment(s storage.Comment) error {
	err := c.db.QueryRow(context.Background(),
		`DELETE FROM comments WHERE id = $1;`, s.ID).Scan()
	if err != nil {
		return err
	}
	return nil
}

func (p *CommentDB) CreateCommentTable() error {
	_, err := p.db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS comments (
                id SERIAL PRIMARY KEY,
                id_news INT,
                content TEXT NOT NULL DEFAULT 'empty',
                pubtime BIGINT NOT NULL DEFAULT extract (epoch from now())
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

// DropCommentTable Удаляет таблицу
func (p *CommentDB) DropCommentTable() error {
	_, err := p.db.Exec(context.Background(), "DROP TABLE IF EXISTS comments;")
	if err != nil {
		return err
	}
	return nil
}

package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"news/pkg/storage"
)

type PostsDB struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, databaseUrl string) (*PostsDB, error) {
	for {
		_, err := pgxpool.Connect(ctx, databaseUrl)
		if err == nil {
			break
		}
	}
	db, err := pgxpool.Connect(ctx, databaseUrl)
	if err != nil {
		return nil, err
	}
	p := PostsDB{
		db: db,
	}
	return &p, nil
}

func (p *PostsDB) GetPosts(n int) ([]storage.Post, error) {
	rows, err := p.db.Query(context.Background(), `SELECT id, title, content, pubDate, link 
       FROM posts LIMIT $1;`, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var postsList []storage.Post
	for rows.Next() {
		var stPost storage.Post
		err = rows.Scan(
			&stPost.ID,
			&stPost.Title,
			&stPost.Content,
			&stPost.PubTime,
			&stPost.Link,
		)
		if err != nil {
			return nil, err
		}
		postsList = append(postsList, stPost)
	}
	return postsList, rows.Err()
}

func (p *PostsDB) AddPost(s storage.Post) error {
	err := p.db.QueryRow(context.Background(),
		`INSERT INTO posts (title, content, pubTime, link) VALUES 
             ($1, $2, $3, $4);`, s.Title, s.Content, s.PubTime, s.Link).Scan()
	return err
}

func (p *PostsDB) PostDetail(id int) (storage.Post, error) {
	row := p.db.QueryRow(context.Background(), `SELECT * FROM posts WHERE id=$1;`, id)
	var post storage.Post
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.PubTime, &post.Link)
	if err != nil {
		return storage.Post{}, err
	}
	return post, nil
}

func (p *PostsDB) PostSearch(pattern string, limit, offset int) ([]storage.Post, storage.Pagination, error) {
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

	rows, err := p.db.Query(context.Background(),
		"SELECT * FROM posts WHERE title ILIKE $1 ORDER BY pubtime DESC LIMIT $2 OFFSET $3;",
		pattern, limit, offset)
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

func (s *PostsDB) PostPage(limit, offset int) ([]storage.Post, error) {
	pagination := storage.Pagination{
		Page:  offset/limit + 1,
		Limit: limit,
	}
	rows, err := s.db.Query(context.Background(), `
	SELECT * FROM posts
	ORDER BY pubtime DESC LIMIT $1 OFFSET $2
	`,
		pagination.Limit, pagination.Page,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []storage.Post
	// итерированное по результату выполнения запроса
	// и сканирование каждой строки в переменную
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
		// добавление переменной в массив результатов
		posts = append(posts, p)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return posts, rows.Err()
}

func (p *PostsDB) CreatePostsTable() error {
	_, err := p.db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS posts (
                id SERIAL PRIMARY KEY,
                title TEXT NOT NULL DEFAULT 'empty',
                content TEXT NOT NULL DEFAULT 'empty',
                pubtime BIGINT NOT NULL DEFAULT extract (epoch from now()),
                link TEXT NOT NULL
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

// DropGonewsTable Удаляет таблицу
func (p *PostsDB) DropPostsTable() error {
	_, err := p.db.Exec(context.Background(), "DROP TABLE IF EXISTS posts;")
	if err != nil {
		return err
	}
	return nil
}

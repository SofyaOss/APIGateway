package db

import (
	"censors/pkg/storage"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

// Store Хранилище данных
type Store struct {
	db *pgxpool.Pool
}

// New Конструктор объекта хранилища
func New(ctx context.Context, constr string) (*Store, error) {

	for {
		_, err := pgxpool.Connect(ctx, constr)
		if err == nil {
			break
		}
	}
	db, err := pgxpool.Connect(ctx, constr)
	if err != nil {
		return nil, err
	}
	s := Store{
		db: db,
	}
	return &s, nil
}

func (p *Store) GetList() ([]storage.Stop, error) {
	rows, err := p.db.Query(context.Background(), "SELECT * FROM ban")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []storage.Stop
	for rows.Next() {
		var c storage.Stop
		err = rows.Scan(&c.ID, &c.BannedWords)
		if err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (p Store) AddList(c storage.Stop) error {
	_, err := p.db.Exec(context.Background(),
		"INSERT INTO ban (ban_words) VALUES ($1);", c.BannedWords)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// CreateStopTable Создает таблицу
func (p *Store) CreateStopTable() error {
	_, err := p.db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS ban (
			id SERIAL PRIMARY KEY,
			ban_words TEXT
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

// DropStopTable Удаляет таблицу
func (p *Store) DropStopTable() error {
	_, err := p.db.Exec(context.Background(), "DROP TABLE IF EXISTS ban;")
	if err != nil {
		return err
	}
	return nil
}

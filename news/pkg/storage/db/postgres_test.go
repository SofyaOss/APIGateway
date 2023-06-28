package db

import (
	"context"
	db "skillfactory/36/pkg/storage"
	"testing"
	"time"
)

//const DBURL = "postgres://postgres:postgres@5432/postsDB"

const DBURL = "user=postgres password=Keks17sql dbname=postsDB sslmode=disable"

func TestNew(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, err := New(ctx, DBURL)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostsDB_AddPost(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	newDB, err := New(ctx, DBURL)
	if err != nil {
		t.Fatal(err)
	}
	newPost := db.Post{Title: "Тест", Content: "Текст для теста", PubTime: 5, Link: "Тестовая ссылка"}
	newDB.AddPost(newPost)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Запись создана")
}

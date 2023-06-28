package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"skillfactory/36/pkg/storage"
	"skillfactory/36/pkg/storage/db"
	"testing"
	"time"
)

const DBURL = "postgres://postgres:postgres@5432/postsDB"

func TestAPI_posts(t *testing.T) {
	dataLength := 10
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	newDB, _ := db.New(ctx, DBURL)
	newDB.AddPost(storage.Post{})
	newAPI := New(newDB)
	req := httptest.NewRequest(http.MethodGet, "/news/10", nil)
	newRec := httptest.NewRecorder()
	newAPI.r.ServeHTTP(newRec, req)
	if !(newRec.Code == http.StatusOK) {
		t.Error("Неверный код:", newRec.Code)
	}
	body, err := ioutil.ReadAll(newRec.Body)
	if err != nil {
		t.Fatal("Ошибка при раскодировании:", err)
	}
	var postList []storage.Post
	err = json.Unmarshal(body, &postList)
	if err != nil {
		t.Fatal("Ошибка при раскодировании:", err)
	}
	if len(postList) != dataLength {
		t.Fatalf("Получено неверное количество записей: %d вместо %d", len(postList), dataLength)
	}
}

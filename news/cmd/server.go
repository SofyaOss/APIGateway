package main

import (
	"context"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"news/pkg/api"
	"news/pkg/middleware"
	"news/pkg/rss"
	"news/pkg/storage"
	"news/pkg/storage/db"
	"os"
	"time"
)

//const DBURL = "postgres://postgres:postgres@127.0.0.1:8081/posts"

type server struct {
	db  storage.Interface
	api *api.API
}

type config struct {
	Period  int
	LinkArr []string
}

type Config struct {
	Censor   Censor
	Comments Comments
	News     News
	Gateway  Gateway
}

type Censor struct {
	AdrPort string
	URLdb   string
}

type Comments struct {
	AdrPort string
	URLdb   string
}

type News struct {
	AdrPort string
	URLdb   string
}

type Gateway struct {
	AdrPort string
}

func NewConf() *Config {
	return &Config{
		Censor: Censor{
			AdrPort: getEnv("CENSOR_PORT", ""),
			URLdb:   getEnv("CENSOR_DB", ""),
		},
		Comments: Comments{
			AdrPort: getEnv("COMMENTS_PORT", ""),
			URLdb:   getEnv("COMMENTS_DB", ""),
		},
		News: News{
			AdrPort: getEnv("NEWS_PORT", ""),
			URLdb:   getEnv("NEWS_DB", ""),
		},
		Gateway: Gateway{
			AdrPort: getEnv("GATEWAY_PORT", ""),
		},
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

const (
	configURL = "./cmd/config.json"
)

// Простая вспомогательная функция для считывания окружения или возврата значения по умолчанию
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func main() {
	chanErr, chanPost := make(chan error), make(chan []storage.Post)
	conf := NewConf()
	dbURL := conf.News.URLdb
	port := conf.News.AdrPort
	portFlag := flag.String("news-port", port, "Порт для сервиса news")
	flag.Parse()
	portNews := *portFlag

	var serv server
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	newDB, err := db.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	db, err := db.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	err = db.DropPostsTable()
	if err != nil {
		log.Println(err)
		return
	}
	// Создание таблицы gonews если не существует
	err = db.CreatePostsTable()
	if err != nil {
		log.Println(err)
		return
	}
	serv.db = newDB
	serv.api = api.New(serv.db)
	go func() {
		err := rss.GetNews(configURL, chanPost, chanErr)
		if err != nil {
			log.Fatal(err)
		}
	}()
	go func() { // обработка новостей
		for allPosts := range chanPost {
			for idx := range allPosts {
				newDB.AddPost(allPosts[idx])
			}
		}
	}()
	go func() { // обработка ошибок
		for err := range chanErr {
			log.Println("Ошибка:", err)
		}
	}()
	serv.api.Router().Use(middleware.Middle)
	log.Print("Запуск сервера на http://127.0.0.1" + portNews)
	err = http.ListenAndServe(portNews, serv.api.Router())
	if err != nil {
		log.Fatal(err)
	}
}

//func main() {
//	chanErr, chanPost := make(chan error), make(chan []storage.Post)
//	conFile, err := ioutil.ReadFile("./config.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//	var config config
//	err = json.Unmarshal(conFile, &config)
//	if err != nil {
//		log.Fatal(err)
//	}
//	rssLinks := rssJSON("config.json", chanErr)
//	for i := range rssLinks.LinkArr {
//		go postParse(rssLinks.LinkArr[i], config.Period, chanErr, chanPost)
//	}
//	var serv server
//	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
//	defer cancel()
//	newDB, err := db.New(ctx, DBURL)
//	if err != nil {
//		log.Fatal(err)
//	}
//	serv.db = newDB
//	go func() { // обработка новостей
//		for allPosts := range chanPost {
//			for idx := range allPosts {
//				newDB.AddPost(allPosts[idx])
//			}
//		}
//	}()
//	go func() { // обработка ошибок
//		for err := range chanErr {
//			log.Println("Ошибка:", err)
//		}
//	}()
//	err = http.ListenAndServe(":80", serv.api.Router())
//	if err != nil {
//		log.Fatal(err)
//	}
//}

//func postParse(link string, dur int, chanErr chan<- error, chanPost chan<- []storage.Post) {
//	for {
//		postList, err := rss.getRssStruct(link)
//		if err != nil {
//			chanErr <- err
//			continue
//		}
//		chanPost <- postList
//		time.Sleep(time.Duration(dur) * time.Minute)
//	}
//}
//
//func rssJSON(file string, chanErr chan<- error) config {
//	jsFile, err := os.Open(file)
//	if err != nil {
//		chanErr <- err
//	}
//	defer jsFile.Close()
//	byteValue, _ := ioutil.ReadAll(jsFile)
//	var linkList config
//	json.Unmarshal(byteValue, &linkList)
//	return linkList
//}

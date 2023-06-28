package main

import (
	"censors/config"
	"censors/pkg/api"
	"censors/pkg/middleware"
	"censors/pkg/storage"
	postgres "censors/pkg/storage/db"
	"context"
	"flag"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func bannedWords() ([]storage.Stop, error) {
	file, err := os.Open("./banned_words.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	line := strings.Split(string(content), "\n")
	var res []storage.Stop
	for _, item := range line {
		trimLine := strings.TrimSpace(item)
		newItem := storage.Stop{BannedWords: trimLine}
		res = append(res, newItem)
	}
	return res, nil
}

type server struct {
	db  storage.Interface
	api *api.API
}

// init вызывается перед main()
func init() {
	// загружает значения из файла .env в систему
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	//chanErr, chanPost := make(chan error), make(chan []storage.Post)
	conf := config.New()
	dbURL := conf.Comments.URLdb
	port := conf.Censor.AdrPort
	portFlag := flag.String("news-port", port, "Порт для сервиса censor")
	flag.Parse()
	portCensor := *portFlag

	var serv server
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	db, err := postgres.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	err = db.DropStopTable()
	if err != nil {
		log.Println(err)
		return
	}
	// Создание таблицы stop если не существует
	err = db.CreateStopTable()
	if err != nil {
		log.Println(err)
		return
	}
	// Получение списка для стоп листа из файла words.txt
	stop, err := bannedWords()
	if err != nil {
		log.Println(err)
	}
	// Добавление в таблицу stop полученного списка
	for _, v := range stop {
		err := db.AddList(v)
		if err != nil {
			log.Println(err)
		}
	}

	// Инициализируем хранилище сервера конкретной БД.
	serv.db = db

	// Создаём объект API и регистрируем обработчики.
	serv.api = api.New(serv.db)

	serv.api.Router().Use(middleware.Middle)

	log.Print("Запуск сервера на http://127.0.0.1" + portCensor)

	err = http.ListenAndServe(portCensor, serv.api.Router())
	if err != nil {
		log.Fatal("Не удалось запустить сервер шлюза. Error:", err)
	}
}

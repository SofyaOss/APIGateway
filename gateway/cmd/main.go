package main

import (
	"flag"
	"gateway/config"
	"gateway/pkg/api"
	"gateway/pkg/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

type server struct {
	api *api.API
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	var serv server

	// Конфигурация
	conf := config.New()
	port := conf.Gateway.AdrPort
	newsPort := conf.News.AdrPort
	censorPort := conf.Censor.AdrPort
	comment := conf.Comments.AdrPort

	portFlag := flag.String("gateway-port", port, "Порт для gateway сервиса")

	portFlagNews := flag.String("news-port", newsPort, "Порт для news сервиса")

	portFlagCensor := flag.String("censor-port", censorPort, "Порт для censor сервиса")

	portFlagComment := flag.String("comments-port", comment, "Порт для comments сервиса")

	flag.Parse()

	portGateway := *portFlag
	portNews := *portFlagNews
	portCensor := *portFlagCensor
	portComment := *portFlagComment

	serv.api = api.New(conf, portNews, portCensor, portComment)
	serv.api.Router().Use(middleware.Middle)

	log.Print("Запуск сервера http://127.0.0.1" + portGateway + "/news")

	err := http.ListenAndServe(portGateway, serv.api.Router())
	if err != nil {
		log.Fatal("Не удалось запустить сервер шлюза. Error:", err)
	}
}

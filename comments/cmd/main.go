package main

import (
	"comments/pkg/api"
	"comments/pkg/storage"
	"github.com/joho/godotenv"
	"log"
)

type server struct {
	db  storage.Interface
	api *api.API
}

func init() {
	// загружает значения из файла .env в систему
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

}

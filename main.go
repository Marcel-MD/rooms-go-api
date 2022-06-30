package main

import (
	"github.com/Marcel-MD/rooms-go-api/handlers"
	"github.com/Marcel-MD/rooms-go-api/models"
)

func main() {
	models.InitDB()
	r := handlers.InitRouter()
	r.Run()
}

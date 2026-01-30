package main

import (
	"log"

	"github.com/AlikhanF2006/Final_project/configs"
	"github.com/AlikhanF2006/Final_project/pkg/db"
)

func main() {
	configs.LoadConfig()
	db.Connect()
	defer db.Close()

	log.Println("server started")
	select {}
}

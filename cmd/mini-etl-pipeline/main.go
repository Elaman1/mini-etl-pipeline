package main

import (
	"log"
	runserver "mini-etl-pipeline/internal/app"
	"os"
)

func main() {
	if err := runserver.RunServer(); err != nil {
		log.Printf("Вышла ошибка %v", err)
		os.Exit(1)
	}
}

package main

import (
	"3xinfobot/internal/config"
	"log"
)

func main() { 
	log.Printf("Starting 3xInfoBot...")
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)		
	}
	log.Printf("Config loaded successfully")
}
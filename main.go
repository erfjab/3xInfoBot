package main

import (
	"3xinfobot/internal/config"
	"3xinfobot/internal/handlers"
	"3xinfobot/internal/panel"
	"log"
	"time"

	"github.com/erfjab/egobot/core"
)

func main() {
	log.Printf("Starting 3xInfoBot...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Config loaded successfully")

	client, err := panel.NewClient(cfg.PanelURL, cfg.PanelPath)
	if err != nil {
		log.Fatalf("Failed to create panel client: %v", err)
	}

	if err = client.Login(cfg.PanelUsername, cfg.PanelPassword, cfg.PanelTwoFactorCode); err != nil {
		log.Fatalf("Panel login failed: %v", err)
	}
	panel.DefaultClient = client
	log.Printf("Panel login successful")

	go func() {
		ticker := time.NewTicker(15 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if err := panel.DefaultClient.Login(cfg.PanelUsername, cfg.PanelPassword, cfg.PanelTwoFactorCode); err != nil {
				log.Printf("Panel session refresh failed: %v", err)
			} else {
				log.Printf("Panel session refreshed")
			}
		}
	}()

	bot := core.NewBot(cfg.TelegramBotToken)

	getBotInfo, err := bot.GetMe()
	if err != nil {
		log.Fatalf("error getting bot info: %v", err)
	}
	log.Printf("Telegram bot initialized: @%s (ID: %d)", getBotInfo.Username, getBotInfo.ID)

	bot.OnCommand("start", handlers.StartHandler)
	bot.OnMessage(handlers.LinkHandler)

	bot.StartPolling(&core.PollingOptions{})
}

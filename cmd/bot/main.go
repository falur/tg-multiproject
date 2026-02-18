package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"tg-multiproject/internal/bot"
	"tg-multiproject/internal/config"
	"tg-multiproject/internal/state"
	"tg-multiproject/internal/storage"
)

func main() {
	cfg := config.Load()

	if cfg.TelegramToken == "" {
		log.Fatal("TELEGRAM_TOKEN is required")
	}

	// Ensure directories exist
	if err := os.MkdirAll(cfg.ProjectsDir, 0o755); err != nil {
		log.Fatalf("create projects dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(cfg.DatabasePath), 0o755); err != nil {
		log.Fatalf("create data dir: %v", err)
	}

	store, err := storage.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("init storage: %v", err)
	}
	defer store.Close()

	sm := state.NewManager()

	b, err := bot.New(cfg, store, sm)
	if err != nil {
		log.Fatalf("init bot: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go b.Start()

	log.Println("Bot started")
	<-ctx.Done()

	log.Println("Shutting down...")
	b.Stop()
}

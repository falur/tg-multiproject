package config

import "os"

type Config struct {
	TelegramToken string
	AllowedUserID int64
	ClaudeBinary  string
	ProjectsDir   string
	DatabasePath  string
}

func Load() *Config {
	uid := envOrDefault("ALLOWED_USER_ID", "0")
	var allowedUID int64
	for _, ch := range uid {
		if ch >= '0' && ch <= '9' {
			allowedUID = allowedUID*10 + int64(ch-'0')
		}
	}

	return &Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		AllowedUserID: allowedUID,
		ClaudeBinary:  envOrDefault("CLAUDE_BINARY", "claude"),
		ProjectsDir:   envOrDefault("PROJECTS_DIR", "./projects"),
		DatabasePath:  envOrDefault("DATABASE_PATH", "./data/bot.db"),
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	NotionToken string
	PageID      string
	UserID      string
	Location    *time.Location
}

func LoadConfig() *Config {
	// 로컬 개발 환경일 때만 .env 파일 로드
	_ = godotenv.Load()
	loc := time.FixedZone("KST", 9*60*60)

	return &Config{
		NotionToken: os.Getenv("NOTION_TOKEN"),
		PageID:      os.Getenv("NOTION_PAGE_ID"),
		UserID:      os.Getenv("NOTION_USER_ID"),
		Location:    loc,
	}
}

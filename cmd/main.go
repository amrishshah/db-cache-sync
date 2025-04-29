package main

import (
	"log"

	"github.com/amrishkshah/db-cache-sync/internal/binlog"
	"github.com/amrishkshah/db-cache-sync/internal/cache"
	"github.com/amrishkshah/db-cache-sync/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	redisClient := cache.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword)

	log.Println("Starting binlog reader...")
	err := binlog.StartBinlogReader(cfg, redisClient)
	if err != nil {
		log.Fatalf("Error starting binlog reader: %v", err)
	}
}

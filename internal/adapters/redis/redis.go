package redis

import (
	"context"
	"log"
	"strings"

	goredis "github.com/redis/go-redis/v9"
)

func NewClient(addrOrURL string) *goredis.Client {
	var rdb *goredis.Client

	if strings.HasPrefix(addrOrURL, "redis://") || strings.HasPrefix(addrOrURL, "rediss://") {
		opts, err := goredis.ParseURL(addrOrURL)
		if err != nil {
			log.Fatalf("redis parse url: %v", err)
		}
		rdb = goredis.NewClient(opts)
	} else {
		rdb = goredis.NewClient(&goredis.Options{Addr: addrOrURL})
	}

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis ping: %v", err)
	}
	log.Println("connected to redis")
	return rdb
}
package libnord

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Timeout       time.Duration
	Prefix        string
	Cache         bool
	CacheTimeout  int
	RedisAddr     string
	RedisPassword string
	Logger        *log.Logger
}

var DefaultConfig = &Config{
	Timeout:       time.Duration(3000 * time.Millisecond),
	Prefix:        "",
	Cache:         false,
	CacheTimeout:  60,
	RedisAddr:     ":6379",
	RedisPassword: "",
	Logger:        log.New(os.Stdout, "[nord] ", 3),
}

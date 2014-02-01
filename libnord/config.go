package libnord

import (
	"time"
)

type Config struct {
	Timeout       time.Duration
	Prefix        string
	Cache         bool
	CacheTimeout  int
	RedisAddr     string
	RedisPassword string
}

var DefaultConfig = &Config{
	Timeout:       time.Duration(3000 * time.Millisecond),
	Prefix:        "",
	Cache:         false,
	CacheTimeout:  60,
	RedisAddr:     ":6379",
	RedisPassword: "",
}

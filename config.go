package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/FZambia/nord/libnord"
)

type Config struct {
	Addr string
	*libnord.Config
}

func parseCommandLineOptions() Config {
	var mainConfig Config
	var config libnord.Config

	portFlag := flag.Int("port", 3000, "HTTP port to listen on")
	addressFlag := flag.String("address", "0.0.0.0", "Interface to bind to (e.g. 127.0.0.1)")
	timeoutFlag := flag.Int("timeout", 3000, "Default request timeout in milliseconds")
	prefixFlag := flag.String("prefix", "", "Url prefix for application handlers (e.g. /count/)")

	cacheFlag := flag.Bool("cache", false, "Use Redis for caching")
	cacheTimeoutFlag := flag.Int("cache-timeout", 60, "Cache timeout in seconds")
	redisPortFlag := flag.Int("redis-port", 6379, "Redis port")
	redisAddressFlag := flag.String("redis-address", "127.0.0.1", "Redis interface (e.g. 127.0.0.1)")
	redisPasswordFlag := flag.String("redis-password", "", "Redis auth password, not used by default")

	flag.Parse()

	mainConfig.Addr = fmt.Sprintf("%s:%d", *addressFlag, *portFlag)

	config.Timeout = time.Duration(*timeoutFlag) * time.Millisecond
	config.Prefix = *prefixFlag

	config.Cache = *cacheFlag
	config.CacheTimeout = *cacheTimeoutFlag
	config.RedisAddr = fmt.Sprintf("%s:%d", *redisAddressFlag, *redisPortFlag)
	config.RedisPassword = *redisPasswordFlag

	mainConfig.Config = &config
	return mainConfig
}

package libnord

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

const cachePrefix = "nord"

func getRedisConnection(config *Config) (redis.Conn, error) {
	if !config.Cache {
		return nil, nil
	}
	conn, err := redis.Dial("tcp", config.RedisAddr)
	if err != nil {
		log.Println("Redis connection error:", err)
		return nil, err
	}
	if len(config.RedisPassword) > 0 {
		if _, err := conn.Do("AUTH", config.RedisPassword); err != nil {
			conn.Close()
			return nil, err
		}
	}
	log.Println("connected with Redis on", config.RedisAddr)
	return conn, err
}

func reconnectToRedis(service *ServiceMap) {
	service.conn.Close()
	conn, err := getRedisConnection(service.config)
	if err != nil {
		log.Println(err)
		return
	}
	service.conn = conn
}

func getCacheKey(providerName, url string) string {
	return fmt.Sprintf("%s:%s:%s", cachePrefix, providerName, url)
}

func getCacheResult(service *ServiceMap, providerName, url string) (interface{}, error) {
	if !service.config.Cache {
		return nil, nil
	}
	reply, err := service.conn.Do("GET", getCacheKey(providerName, url))
	if err != nil {
		reconnectToRedis(service)
		return nil, err
	}
	if reply == nil {
		return nil, nil
	}
	reply_bytes, err := redis.Bytes(reply, err)
	if err != nil {
		return nil, err
	}
	var f interface{}
	err = json.Unmarshal(reply_bytes, &f)
	return f, err
}

func setCacheResult(service *ServiceMap, providerName, url string, resp *ProviderResponse) error {
	if !service.config.Cache {
		return nil
	}
	key := getCacheKey(providerName, url)
	value, _ := json.Marshal(resp)
	service.conn.Send("MULTI")
	service.conn.Send("SET", key, value)
	service.conn.Send("EXPIRE", key, service.config.CacheTimeout)
	_, err := service.conn.Do("EXEC")
	return err
}

package redis

import (
	"fmt"
	"os"
	"strconv"

	rs "github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

func GetInstance(log *logrus.Logger) *rs.Client {
	log.Infoln("Connecting to Redis...")
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	addr := fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	redisConn := rs.NewClient(&rs.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
		Network:  os.Getenv("REDIS_NETWORK"),
	})
	_, err := redisConn.Ping().Result()
	if err != nil {
		log.Fatalf("Redis connection failed: %v\n", err)
	}
	log.Infoln("Redis started at: " + addr)
	return redisConn
}

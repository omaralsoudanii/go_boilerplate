package redis

import (
	rs "github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

func GetInstance(log *logrus.Logger) *rs.Client {
	redisConn := rs.NewClient(&rs.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := redisConn.Ping().Result()
	if err != nil {
		log.Fatalf("Redis connection failed: %v\n", err)
	}
	return redisConn
}

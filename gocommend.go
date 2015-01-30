package gocommend

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

var (
	redisClient redis.Conn
	err         error
)

func init() {
	redisClient, err = redis.Dial("tcp", localRedisURL+":"+localRedisPort)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func Redistest() {
	redisClient.Do("SET", "aaa", 123)
	a, err := redis.Int(redisClient.Do("GET", "aaa"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(a)
}

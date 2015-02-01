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
	if localStartup == true {
		redisClient, err = redis.Dial("tcp", localRedisURL+":"+localRedisPort)
	} else {
		redisClient, err = redis.Dial("tcp", remoteRedisURL+":"+remoteRedisPort)
	}

	if err != nil {
		log.Println(err.Error())
		return
	}
}

func Redistest() {
	//redisClient.Do("SET", "aaa", 123)
	//a, err := redis.Int(redisClient.Do("GET", "aaa"))
	a, err := redisClient.Do("SUNION", "movie:item44:itemLiked", "movie:item44:itemDisliked")
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(a)
}

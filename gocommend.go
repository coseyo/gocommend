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
	if LOCAL_STARTUP == true {
		redisClient, err = redis.Dial("tcp", LOCAL_REDIS_HOST+":"+LOCAL_REDIS_PORT)
	} else {
		redisClient, err = redis.Dial("tcp", REMOTE_REDIS_HOST+":"+REMOTE_REDIS_PORT)
	}
	if err != nil {
		log.Println(err.Error())
		return
	}
}

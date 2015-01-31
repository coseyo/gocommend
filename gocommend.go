package gocommend

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

var (
	redisClient redis.Conn
	err         error
)

type Option struct {
	collection string
}

type CollectionSet struct {
	liked          string
	disliked       string
	userSimilarity string
	itemSimilarity string
	temp           string
	tempDiff       string
	userRecommend  string
}

func (o *Option) SetCollection(name string) {
	o.collection = name
}

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
	redisClient.Do("SET", "aaa", 123)
	a, err := redis.Int(redisClient.Do("GET", "aaa"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(a)
}

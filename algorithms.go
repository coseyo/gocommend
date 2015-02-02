package gocommend

import (
	"log"
	"math"

	"github.com/garyburd/redigo/redis"
)

type algorithms struct {
	cSet collectionSet
}

func (this *algorithms) updateWilsonScore(itemId string) error {
	var (
		total int
		pOS   float64
		score float64 = 0.0
	)

	resultLike, _ := redis.Int(redisClient.Do("SCARD", this.cSet.itemLiked(itemId)))
	resultDislike, _ := redis.Int(redisClient.Do("SCARD", this.cSet.itemDisliked(itemId)))

	total = resultLike + resultDislike
	if total > 0 {
		pOS = float64(resultLike) / float64(total)
		score = this.willsonScore(total, pOS)
	}

	_, err := redisClient.Do("ZADD", this.cSet.scoreRank, score, itemId)

	return err
}

func (this *algorithms) willsonScore(total int, pOS float64) float64 {

	var z float64 = 1.96

	n := float64(total)

	return math.Abs((pOS + z*z/(2*n) - z*math.Sqrt(pOS*(1-pOS)+z*z/(4*n))) / (1 + z*z/n))
}

func (this *algorithms) similaritySum(simSet string, compSet string) float64 {
	var similarSum float64 = 0.0
	userIds, err := redis.Values(redisClient.Do("SMEMBERS", compSet))
	log.Println(compSet)
	for _, rs := range userIds {
		userId, _ := redis.String(rs, err)
		log.Println(userIds)
		score, _ := redis.Float64(redisClient.Do("ZSCORE", simSet, userId))
		similarSum += score
	}
	return similarSum
}

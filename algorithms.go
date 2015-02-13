package gocommend

import "github.com/garyburd/redigo/redis"

// algorithm type's parent
type algorithms struct {
	cSet collectionSet
}

// 2 set's similarity
func (this *algorithms) similaritySum(simSet string, compSet string) float64 {
	var similarSum float64 = 0.0
	userIds, err := redis.Values(redisClient.Do("SMEMBERS", compSet))
	for _, rs := range userIds {
		userId, _ := redis.String(rs, err)
		score, _ := redis.Float64(redisClient.Do("ZSCORE", simSet, userId))
		similarSum += score
	}
	return similarSum
}

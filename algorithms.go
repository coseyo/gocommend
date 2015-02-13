package gocommend

import "github.com/garyburd/redigo/redis"

// algorithm type's parent
type algorithms struct {
	cSet collectionSet
}

func (this *algorithms) TrimRecommendItem(userId string) error {
	count, err := redis.Int(redisClient.Do("ZCARD", this.cSet.recommendedItem(userId)))
	if err != nil {
		return err
	}
	if count > MAX_RECOMMEND_ITEM {
		redisClient.Do("ZREMRANGEBYRANK", this.cSet.recommendedItem(userId), 0, (count - MAX_RECOMMEND_ITEM))
	}
	return nil
}

func (this *algorithms) TrimUserSimilarity(userId string) error {
	count, err := redis.Int(redisClient.Do("ZCARD", this.cSet.userSimilarity(userId)))
	if err != nil {
		return err
	}
	if count > MAX_SIMILARITY_USER {
		redisClient.Do("ZREMRANGEBYRANK", this.cSet.userSimilarity(userId), 0, (count - MAX_SIMILARITY_USER))
	}
	return nil
}

func (this *algorithms) TrimItemSimilarity(itemId string) error {
	count, err := redis.Int(redisClient.Do("ZCARD", this.cSet.itemSimilarity(itemId)))
	if err != nil {
		return err
	}
	if count > MAX_SIMILARITY_ITEM {
		redisClient.Do("ZREMRANGEBYRANK", this.cSet.itemSimilarity(itemId), 0, (count - MAX_SIMILARITY_ITEM))
	}
	return nil
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

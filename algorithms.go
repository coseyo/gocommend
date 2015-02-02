package gocommend

import (
	"log"
	"math"

	"github.com/garyburd/redigo/redis"
)

type algorithms struct {
	cSet collectionSet
}

func (algo *algorithms) updateSimilarityFor(userId string) error {
	if ratedItemSet, err := redis.Values(redisClient.Do("SUNION", algo.cSet.userLiked(userId), algo.cSet.userDisliked(userId))); err != nil {
		return err
	}

	if len(ratedItemSet) == 0 {
		return nil
	}

	itemLikeDislikeKeys := []string{}
	for _, rs := range ratedItemSet {
		itemId, _ := redis.String(rs, err)
		itemLikeDislikeKeys = append(itemLikeDislikeKeys, algo.cSet.itemLiked(itemId))
		itemLikeDislikeKeys = append(itemLikeDislikeKeys, algo.cSet.itemDisliked(itemId))
	}

	otherUserIdsWhoRated, err := redis.Values(redisClient.Do("SUNION", redis.Args{}.AddFlat(itemLikeDislikeKeys)...))

	if err != nil {
		log.Panicln("error sunion 2")
		return err
	}

	for _, rs := range otherUserIdsWhoRated {
		otherUserId, _ := redis.String(rs, err)

		if len(otherUserIdsWhoRated) == 1 || userId == otherUserId {
			continue
		}

		log.Println(otherUserId)

		score := algo.jaccardCoefficient(userId, otherUserId)

		redisClient.Do("ZADD", algo.cSet.userSimilarity(userId), score, otherUserId)

		log.Println(score)
	}

	return err
}

func (algo *algorithms) jaccardCoefficient(userId1 string, userId2 string) float64 {
	var (
		similarity   int = 0
		rateInCommon int = 0
	)

	resultBothLike, _ := redis.Values(redisClient.Do("SINTER", algo.cSet.userLiked(userId1), algo.cSet.userLiked(userId2)))
	resultBothDislike, _ := redis.Values(redisClient.Do("SINTER", algo.cSet.userDisliked(userId1), algo.cSet.userDisliked(userId2)))
	resultUser1LikeUser2Dislike, _ := redis.Values(redisClient.Do("SINTER", algo.cSet.userLiked(userId1), algo.cSet.userDisliked(userId2)))
	resultUser1DislikeUser2Like, _ := redis.Values(redisClient.Do("SINTER", algo.cSet.userDisliked(userId1), algo.cSet.userLiked(userId2)))

	len1 := len(resultBothLike)
	len2 := len(resultBothDislike)
	len3 := len(resultUser1LikeUser2Dislike)
	len4 := len(resultUser1DislikeUser2Like)

	similarity = len1 + len2 - len3 - len4
	rateInCommon = len1 + len2 + len3 + len4

	return float64(similarity) / float64(rateInCommon)
}

func (algo *algorithms) updateWilsonScore(itemId string) error {
	var (
		total int
		pOS   float64
		score float64 = 0.0
	)

	resultLike, _ := redis.Int(redisClient.Do("SCARD", algo.cSet.itemLiked(itemId)))
	resultDislike, _ := redis.Int(redisClient.Do("SCARD", algo.cSet.itemDisliked(itemId)))

	total = resultLike + resultDislike
	if total > 0 {
		pOS = float64(resultLike) / float64(total)
		score = algo.willsonScore(total, pOS)
	}

	_, err := redisClient.Do("ZADD", algo.cSet.scoreRank, score, itemId)

	return err
}

func (algo *algorithms) willsonScore(total int, pOS float64) float64 {

	var z float64 = 1.96

	n := float64(total)

	return math.Abs((pOS + z*z/(2*n) - z*math.Sqrt(pOS*(1-pOS)+z*z/(4*n))) / (1 + z*z/n))
}

func (algo *algorithms) updateRecommendationFor(userId string) error {

	mostSimilarUserIds, err := redis.Values(redisClient.Do("ZREVRANGE", algo.cSet.userSimilarity(userId), 0, MAX_NEIGHBORS-1))

	if len(mostSimilarUserIds) == 0 {
		return err
	}

	for _, rs := range mostSimilarUserIds {
		similarUserId, _ := redis.String(rs, err)
		redisClient.Do("SUNIONSTORE", algo.cSet.userTemp(userId), algo.cSet.userLiked(similarUserId))
	}

	diffItemIds, err := redis.Values(redisClient.Do("SDIFF", algo.cSet.userLiked(userId), algo.cSet.userDisliked(userId)))

	for _, rs := range diffItemIds {
		diffItemId, _ := redis.String(rs, err)
		score := algo.predictFor(userId, diffItemId)
		redisClient.Do("ZADD", algo.cSet.recommendedItem(userId), score, userId)
	}

	recNum, err := redis.Int(redisClient.Do("ZCARD", algo.cSet.recommendedItem(userId)))

	log.Println("recNum: ", recNum)

	if recNum > MAX_RECOMMEND_ITEM {
		redisClient.Do("ZREMRANGEBYRANK", algo.cSet.recommendedItem(userId), MAX_RECOMMEND_ITEM, -1)
	}
	redisClient.Do("DEL", algo.cSet.userTemp(userId))
	return err
}

func (algo *algorithms) predictFor(userId string, itemId string) float64 {

	result1 := algo.similaritySum(algo.cSet.userSimilarity(userId), algo.cSet.itemLiked(itemId))

	result2 := algo.similaritySum(algo.cSet.userSimilarity(userId), algo.cSet.itemDisliked(itemId))

	log.Println("predict userId:", userId)
	log.Println("predict itemId:", itemId)
	log.Println("predict result 1:", result1)
	log.Println("predict result 2:", result2)

	sum := result1 - result2

	itemLikedCount, _ := redis.Int(redisClient.Do("SCARD", algo.cSet.itemLiked(itemId)))

	itemDislikedCount, _ := redis.Int(redisClient.Do("SCARD", algo.cSet.itemLiked(itemId)))

	return float64(sum) / float64(itemLikedCount+itemDislikedCount)
}

func (algo *algorithms) similaritySum(simSet string, compSet string) float64 {
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

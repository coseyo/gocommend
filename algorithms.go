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
	ratedItemSet, err := redis.Values(redisClient.Do("SUNION", algo.cSet.userLiked(userId), algo.cSet.userDisliked(userId)))

	itemLikeDislikeKeys := []string{}

	if len(ratedItemSet) > 0 {
		for _, rs := range ratedItemSet {
			itemId, _ := redis.String(rs, err)
			itemLikeDislikeKeys = append(itemLikeDislikeKeys, algo.cSet.itemLiked(itemId))
			itemLikeDislikeKeys = append(itemLikeDislikeKeys, algo.cSet.itemDisliked(itemId))
		}
	}

	otherUserIdsWhoRated, err := redis.Values(redisClient.Do("SUNION", redis.Args{}.AddFlat(itemLikeDislikeKeys)...))

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
		similarity        int     = 0
		rateInCommon      int     = 0
		finalJaccardScore float64 = 0.0
	)

	resultBothLike, _ := redis.Values(redisClient.Do("SINTER", algo.cSet.userLiked(userId1), algo.cSet.userLiked(userId2)))
	resultBothDislike, _ := redis.Values(redisClient.Do("SINTER", algo.cSet.userDisliked(userId1), algo.cSet.userDisliked(userId2)))
	resultUser1LikeUser2Dislike, _ := redis.Values(redisClient.Do("SINTER", algo.cSet.userLiked(userId1), algo.cSet.userDisliked(userId2)))
	resultUser1DislikeUser2Like, _ := redis.Values(redisClient.Do("SINTER", algo.cSet.userDisliked(userId1), algo.cSet.userLiked(userId2)))

	len1 := len(resultBothLike)
	len2 := len(resultBothDislike)
	len3 := len(resultUser1LikeUser2Dislike)
	len4 := len(resultUser1DislikeUser2Like)

	log.Println(len1, len2, len3, len4)

	similarity = len1 + len2 - len3 - len4
	rateInCommon = len1 + len2 + len3 + len4
	finalJaccardScore = float64(similarity) / float64(rateInCommon)

	return finalJaccardScore
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

//func (algo *glgorithums)

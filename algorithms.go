package gocommend

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

type algo struct {
	cSet collectionSet
}

func (algo *algo) updateSimilarityFor(userId string) error {
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

		log.Println(score)
	}

	return nil
}

func (algo *algo) jaccardCoefficient(userId1 string, userId2 string) float32 {
	var (
		similarity        int     = 0
		rateInCommon      int     = 0
		finalJaccardScore float32 = 0.0
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

	rateInCommon = len1 + len2 + len3 - len4

	finalJaccardScore = float32(similarity) / float32(rateInCommon)

	return finalJaccardScore
}

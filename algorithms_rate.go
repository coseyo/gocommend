package gocommend

import (
	"math"

	"github.com/garyburd/redigo/redis"
)

// rate type
// Use this type when we collet both like and dislike data,
type algorithmsRate struct {
	algorithms
}

func (this *algorithmsRate) updateSimilarityFor(userId string) error {
	ratedItemSet, err := redis.Values(redisClient.Do("SUNION", this.cSet.userLiked(userId), this.cSet.userDisliked(userId)))

	if err != nil {
		return err
	}

	if len(ratedItemSet) == 0 {
		return nil
	}

	itemLikeDislikeKeys := []string{}
	for _, rs := range ratedItemSet {
		itemId, _ := redis.String(rs, err)
		itemLikeDislikeKeys = append(itemLikeDislikeKeys, this.cSet.itemLiked(itemId))
		itemLikeDislikeKeys = append(itemLikeDislikeKeys, this.cSet.itemDisliked(itemId))
	}

	otherUserIdsWhoRated, err := redis.Values(redisClient.Do("SUNION", redis.Args{}.AddFlat(itemLikeDislikeKeys)...))

	if err != nil {
		return err
	}
	if len(otherUserIdsWhoRated) == 1 {
		return nil
	}

	for _, rs := range otherUserIdsWhoRated {
		otherUserId, _ := redis.String(rs, err)
		if userId == otherUserId {
			continue
		}

		score := this.jaccardCoefficient(userId, otherUserId)
		redisClient.Do("ZADD", this.cSet.userSimilarity(userId), score, otherUserId)
	}

	return err
}

func (this *algorithmsRate) jaccardCoefficient(userId1 string, userId2 string) float64 {
	var (
		similarity   int = 0
		rateInCommon int = 0
	)

	resultBothLike, _ := redis.Values(redisClient.Do("SINTER", this.cSet.userLiked(userId1), this.cSet.userLiked(userId2)))
	resultBothDislike, _ := redis.Values(redisClient.Do("SINTER", this.cSet.userDisliked(userId1), this.cSet.userDisliked(userId2)))
	resultUser1LikeUser2Dislike, _ := redis.Values(redisClient.Do("SINTER", this.cSet.userLiked(userId1), this.cSet.userDisliked(userId2)))
	resultUser1DislikeUser2Like, _ := redis.Values(redisClient.Do("SINTER", this.cSet.userDisliked(userId1), this.cSet.userLiked(userId2)))

	len1 := len(resultBothLike)
	len2 := len(resultBothDislike)
	len3 := len(resultUser1LikeUser2Dislike)
	len4 := len(resultUser1DislikeUser2Like)

	similarity = len1 + len2 - len3 - len4
	rateInCommon = len1 + len2 + len3 + len4
	return float64(similarity) / float64(rateInCommon)
}

func (this *algorithmsRate) updateRecommendationFor(userId string) error {

	mostSimilarUserIds, err := redis.Values(redisClient.Do("ZREVRANGE", this.cSet.userSimilarity(userId), 0, MAX_NEIGHBORS-1))

	if len(mostSimilarUserIds) == 0 {
		return err
	}

	for _, rs := range mostSimilarUserIds {
		similarUserId, _ := redis.String(rs, err)
		redisClient.Do("SUNIONSTORE", this.cSet.userTemp(userId), this.cSet.userLiked(similarUserId))
	}

	diffItemIds, err := redis.Values(redisClient.Do("SDIFF", this.cSet.userTemp(userId), this.cSet.userLiked(userId), this.cSet.userDisliked(userId)))
	for _, rs := range diffItemIds {
		diffItemId, _ := redis.String(rs, err)
		score := this.predictFor(userId, diffItemId)
		redisClient.Do("ZADD", this.cSet.recommendedItem(userId), score, diffItemId)
	}

	redisClient.Do("DEL", this.cSet.userTemp(userId))
	return err
}

func (this *algorithmsRate) predictFor(userId string, itemId string) float64 {

	result1 := this.similaritySum(this.cSet.userSimilarity(userId), this.cSet.itemLiked(itemId))

	result2 := this.similaritySum(this.cSet.userSimilarity(userId), this.cSet.itemDisliked(itemId))

	sum := result1 - result2

	itemLikedCount, _ := redis.Int(redisClient.Do("SCARD", this.cSet.itemLiked(itemId)))

	itemDislikedCount, _ := redis.Int(redisClient.Do("SCARD", this.cSet.itemLiked(itemId)))

	return float64(sum) / float64(itemLikedCount+itemDislikedCount)
}

// update socre
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

// willson score
func (this *algorithms) willsonScore(total int, pOS float64) float64 {

	// 95%
	var z float64 = 1.96

	n := float64(total)

	return math.Abs((pOS + z*z/(2*n) - z*math.Sqrt(pOS*(1-pOS)+z*z/(4*n))) / (1 + z*z/n))
}

func (this *algorithmsRate) updateAllData() error {
	userIds, err := redis.Values(redisClient.Do("SMEMBERS", this.cSet.allUser))
	for _, rs := range userIds {
		userId, _ := redis.String(rs, err)
		err = this.updateData(userId, "")
		if err != nil {
			break
		}
	}
	return err
}

func (this *algorithmsRate) updateData(userId string, itemId string) error {

	if err := this.updateSimilarityFor(userId); err != nil {
		return err
	}
	if err := this.updateRecommendationFor(userId); err != nil {
		return err
	}

	if itemId == "" {
		ratedItemSet, err := redis.Values(redisClient.Do("SMEMBERS", this.cSet.userLiked(userId)))
		for _, rs := range ratedItemSet {
			ratedItemId, _ := redis.String(rs, err)
			this.updateWilsonScore(ratedItemId)
		}
	} else {
		if err := this.updateWilsonScore(itemId); err != nil {
			return err
		}
	}
	return nil
}

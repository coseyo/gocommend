package gocommend

import "github.com/garyburd/redigo/redis"

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

	for _, rs := range otherUserIdsWhoRated {
		otherUserId, _ := redis.String(rs, err)
		if len(otherUserIdsWhoRated) == 1 || userId == otherUserId {
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

	recNum, err := redis.Int(redisClient.Do("ZCARD", this.cSet.recommendedItem(userId)))

	if recNum > MAX_RECOMMEND_ITEM {
		redisClient.Do("ZREMRANGEBYRANK", this.cSet.recommendedItem(userId), MAX_RECOMMEND_ITEM, -1)
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

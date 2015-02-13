package gocommend

import "github.com/garyburd/redigo/redis"

// poll type
// we use this type when we don't collect users's dislike data.
type algorithmsPoll struct {
	algorithms
}

func (this *algorithmsPoll) updateUserSimilarity(userId string) error {
	ratedItemSet, err := redis.Values(redisClient.Do("SMEMBERS", this.cSet.userLiked(userId)))

	if err != nil {
		return err
	}

	if len(ratedItemSet) == 0 {
		return nil
	}

	itemKeys := []string{}
	for _, rs := range ratedItemSet {
		itemId, _ := redis.String(rs, err)
		itemKeys = append(itemKeys, this.cSet.itemLiked(itemId))
	}

	otherUserIdsWhoRated, err := redis.Values(redisClient.Do("SUNION", redis.Args{}.AddFlat(itemKeys)...))

	if err != nil {
		return err
	}

	for _, rs := range otherUserIdsWhoRated {
		otherUserId, _ := redis.String(rs, err)
		if len(otherUserIdsWhoRated) == 1 || userId == otherUserId {
			continue
		}

		score := this.jaccardCoefficient(this.cSet.userLiked(userId), this.cSet.userLiked(otherUserId))
		redisClient.Do("ZADD", this.cSet.userSimilarity(userId), score, otherUserId)
	}

	return err
}

func (this *algorithmsPoll) updateItemSimilarity(itemId string) error {
	ratedUserSet, err := redis.Values(redisClient.Do("SMEMBERS", this.cSet.itemLiked(itemId)))

	if err != nil {
		return err
	}

	if len(ratedUserSet) == 0 {
		return nil
	}

	userKeys := []string{}
	for _, rs := range ratedUserSet {
		userId, _ := redis.String(rs, err)
		userKeys = append(userKeys, this.cSet.userLiked(userId))
	}

	otherItemIdsBeingRated, err := redis.Values(redisClient.Do("SUNION", redis.Args{}.AddFlat(userKeys)...))

	if err != nil {
		return err
	}
	if len(otherItemIdsBeingRated) == 1 {
		return nil
	}

	for _, rs := range otherItemIdsBeingRated {
		otherItemId, _ := redis.String(rs, err)
		if itemId == otherItemId {
			continue
		}

		score := this.jaccardCoefficient(this.cSet.itemLiked(itemId), this.cSet.itemLiked(otherItemId))
		redisClient.Do("ZADD", this.cSet.itemSimilarity(itemId), score, otherItemId)
	}

	return err
}

// calculate 2 sets's similarity
func (this *algorithmsPoll) jaccardCoefficient(set1 string, set2 string) float64 {
	var (
		interset int = 0
		unionset int = 0
	)

	resultInter, _ := redis.Values(redisClient.Do("SINTER", set1, set2))
	len1 := len(resultInter)

	len2, _ := redis.Int(redisClient.Do("SCARD", set1))
	len3, _ := redis.Int(redisClient.Do("SCARD", set2))

	interset = len1
	unionset = len2 + len3 - len1
	return float64(interset) / float64(unionset)
}

func (this *algorithmsPoll) updateRecommendationFor(userId string) error {

	mostSimilarUserIds, err := redis.Values(redisClient.Do("ZREVRANGE", this.cSet.userSimilarity(userId), 0, MAX_NEIGHBORS-1))

	if len(mostSimilarUserIds) == 0 {
		return err
	}
	tempSet := this.cSet.userTemp(userId)
	recommendedSet := this.cSet.recommendedItem(userId)

	for _, rs := range mostSimilarUserIds {
		similarUserId, _ := redis.String(rs, err)
		redisClient.Do("SUNIONSTORE", tempSet, this.cSet.userLiked(similarUserId))
	}
	diffItemIds, err := redis.Values(redisClient.Do("SDIFF", tempSet, this.cSet.userLiked(userId)))

	for _, rs := range diffItemIds {
		diffItemId, _ := redis.String(rs, err)
		score := this.predictFor(userId, diffItemId)
		redisClient.Do("ZADD", recommendedSet, score, diffItemId)
	}

	recNum, err := redis.Int(redisClient.Do("ZCARD", this.cSet.recommendedItem(userId)))

	if recNum > MAX_RECOMMEND_ITEM {
		redisClient.Do("ZREMRANGEBYRANK", this.cSet.recommendedItem(userId), MAX_RECOMMEND_ITEM, -1)
	}
	redisClient.Do("DEL", this.cSet.userTemp(userId))
	return err
}

func (this *algorithmsPoll) predictFor(userId string, itemId string) float64 {

	result1 := this.similaritySum(this.cSet.userSimilarity(userId), this.cSet.itemLiked(itemId))

	itemLikedCount, _ := redis.Int(redisClient.Do("SCARD", this.cSet.itemLiked(itemId)))

	return float64(result1) / float64(itemLikedCount)
}

package gocommend

import "github.com/garyburd/redigo/redis"

// input, now support two type of algo
type Input struct {
	cSet collectionSet
}

// init cSet
func (this *Input) Init(collection string) error {
	if collection == "" {
		return gocommendError{emptyCollection}
	}
	this.cSet = collectionSet{}
	this.cSet.init(collection)
	return nil
}

// import rate type data
func (this *Input) ImportRate(userId string, itemId string, rate int) error {
	if rate > 0 {
		if err := like(&this.cSet, userId, itemId); err != nil {
			return err
		}
	} else {
		if err := dislike(&this.cSet, userId, itemId); err != nil {
			return err
		}
	}
	if err := this.UpdateRate(userId, itemId); err != nil {
		return err
	}

	return nil
}

// import poll type data
func (this *Input) ImportPoll(userId string, itemId string) error {

	if err := like(&this.cSet, userId, itemId); err != nil {
		return err
	}
	if err := this.UpdatePoll(userId, itemId); err != nil {
		return err
	}
	return nil
}

// update rate data
func (this *Input) UpdateRate(userId string, itemId string) error {

	algo := algorithmsRate{}
	algo.cSet = this.cSet

	if userId != "" {
		if err := algo.updateSimilarityFor(userId); err != nil {
			return err
		}
		if err := algo.updateRecommendationFor(userId); err != nil {
			return err
		}
	}

	if itemId == "" {
		ratedItemSet, err := redis.Values(redisClient.Do("SMEMBERS", algo.cSet.userLiked(userId)))
		for _, rs := range ratedItemSet {
			ratedItemId, _ := redis.String(rs, err)
			algo.updateWilsonScore(ratedItemId)
		}
	} else {
		if err := algo.updateWilsonScore(itemId); err != nil {
			return err
		}
	}

	return nil
}

// update poll data
func (this *Input) UpdatePoll(userId string, itemId string) error {

	algo := algorithmsPoll{}
	algo.cSet = this.cSet

	if err := algo.updateUserSimilarity(userId); err != nil {
		return err
	}
	if err := algo.updateRecommendationFor(userId); err != nil {
		return err
	}

	if itemId == "" {
		ratedItemSet, err := redis.Values(redisClient.Do("SMEMBERS", algo.cSet.userLiked(userId)))
		for _, rs := range ratedItemSet {
			ratedItemId, _ := redis.String(rs, err)
			algo.updateItemSimilarity(ratedItemId)
		}
	} else {
		if err := algo.updateItemSimilarity(itemId); err != nil {
			return err
		}
	}

	return nil
}

// import original data
func like(cSet *collectionSet, userId string, itemId string) error {
	var (
		rs  interface{}
		err error
	)
	if rs, err = redisClient.Do("SISMEMBER", cSet.itemLiked(itemId), userId); err != nil {
		return err
	}
	if sis, _ := rs.(int); sis == 0 {
		redisClient.Do("ZINCRBY", cSet.mostLiked, 1, itemId)
	}
	if _, err = redisClient.Do("SADD", cSet.allUser, userId); err != nil {
		return err
	}
	if _, err = redisClient.Do("SADD", cSet.userLiked(userId), itemId); err != nil {
		return err
	}
	if _, err = redisClient.Do("SADD", cSet.itemLiked(itemId), userId); err != nil {
		return err
	}
	if _, err = redisClient.Do("ZREM", cSet.recommendedItem(userId), itemId); err != nil {
		return err
	}

	return nil
}

// import original data
func dislike(cSet *collectionSet, userId string, itemId string) error {
	var (
		rs  interface{}
		err error
	)
	if rs, err = redisClient.Do("SISMEMBER", cSet.itemDisliked(itemId), userId); err != nil {
		return err
	}
	if sis, _ := rs.(int); sis == 0 {
		redisClient.Do("ZINCRBY", cSet.mostDisliked, 1, itemId)
	}
	if _, err = redisClient.Do("SADD", cSet.allUser, userId); err != nil {
		return err
	}
	if _, err = redisClient.Do("SADD", cSet.allUser, itemId); err != nil {
		return err
	}
	if _, err = redisClient.Do("SADD", cSet.userDisliked(userId), itemId); err != nil {
		return err
	}
	if _, err = redisClient.Do("SADD", cSet.itemDisliked(itemId), userId); err != nil {
		return err
	}
	if _, err = redisClient.Do("ZREM", cSet.recommendedItem(userId), itemId); err != nil {
		return err
	}

	return nil
}

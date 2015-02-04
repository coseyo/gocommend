package gocommend

import "log"

type Input struct {
	Collection string
	UserId     string
	ItemId     string
	Rate       int
}

func (this *Input) ImportRate() error {
	cSet := collectionSet{}
	cSet.init(this.Collection)

	if this.Rate > 0 {
		log.Println("like")
		if err := like(&cSet, this.UserId, this.ItemId); err != nil {
			return err
		}
	} else {
		log.Println("dislike")
		if err := dislike(&cSet, this.UserId, this.ItemId); err != nil {
			return err
		}
	}
	if err := this.UpdateRate(); err != nil {
		return err
	}

	return nil
}

func (this *Input) ImportPoll() error {
	cSet := collectionSet{}
	cSet.init(this.Collection)
	if err := like(&cSet, this.UserId, this.ItemId); err != nil {
		return err
	}
	if err := this.UpdatePoll(); err != nil {
		return err
	}
	return nil
}

func (this *Input) UpdateRate() error {

	if this.Collection == "" {
		return gocommendError{emptyCollection}
	}

	algo := algorithmsRate{}
	algo.cSet.init(this.Collection)

	// update specific user's sets
	if this.UserId != "" {
		if err := algo.updateSimilarityFor(this.UserId); err != nil {
			return err
		}
		if err := algo.updateRecommendationFor(this.UserId); err != nil {
			return err
		}
	}
	if this.ItemId != "" {
		if err := algo.updateWilsonScore(this.ItemId); err != nil {
			return err
		}
	}

	return nil
}

func (this *Input) UpdatePoll() error {

	if this.Collection == "" {
		return gocommendError{emptyCollection}
	}

	algo := algorithmsPoll{}
	algo.cSet.init(this.Collection)

	// update specific user's sets
	if this.UserId != "" {
		if err := algo.updateSimilarityFor(this.UserId); err != nil {
			return err
		}
		if err := algo.updateRecommendationFor(this.UserId); err != nil {
			return err
		}
	}
	if this.ItemId != "" {
		if err := algo.updateWilsonScore(this.ItemId); err != nil {
			return err
		}
	}

	return nil
}

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

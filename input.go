package gocommend

import "log"

type Input struct {
	Collection string
	UserId     string
	ItemId     string
	Rate       int
}

func Import(i *Input) error {
	log.Println(i)
	cSet, err := initSet(i)
	if err != nil {
		log.Println(err)
	}

	if i.Rate > 0 {
		log.Println("like")
		if err := like(cSet, i.UserId, i.ItemId); err != nil {
			log.Println(err)
			return err
		}
	} else {
		log.Println("dislike")
		if err := dislike(cSet, i.UserId, i.ItemId); err != nil {
			log.Println(err)
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
	if rs, err = redisClient.Do("SISMEMBER", cSet.itemLiked, userId); err != nil {
		return err
	}
	if sis, _ := rs.(int); sis == 0 {
		redisClient.Do("ZINCRBY", cSet.mostLiked, 1, itemId)
	}
	if _, err = redisClient.Do("SADD", cSet.userLiked, itemId); err != nil {
		return err
	}
	if _, err = redisClient.Do("SADD", cSet.itemLiked, userId); err != nil {
		return err
	}

	return nil
}

func dislike(cSet *collectionSet, userId string, itemId string) error {
	var (
		rs  interface{}
		err error
	)
	if rs, err = redisClient.Do("SISMEMBER", cSet.itemDisliked, userId); err != nil {
		return err
	}
	if sis, _ := rs.(int); sis == 0 {
		redisClient.Do("ZINCRBY", cSet.mostDisliked, 1, itemId)
	}
	if _, err = redisClient.Do("SADD", cSet.userDisliked, itemId); err != nil {
		return err
	}
	if _, err = redisClient.Do("SADD", cSet.itemDisliked, userId); err != nil {
		return err
	}

	return nil
}

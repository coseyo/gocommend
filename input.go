package gocommend

import (
	"log"
)

type Input struct {
	Collection string
	UserId     string
	ItemId     string
	Rate       int
}

func Import(i *Input) {
	cSet, err := initSet(i)
	if err != nil {
		log.Println(err)
	}
	log.Println(cSet)

	if i.Rate > 0 {
		if err := like(cSet, i.UserId, i.ItemId); err != nil {
			log.Println(err)
		}
	}
}

func like(cSet *collectionSet, userId string, itemId string) error {
	//if rs, err := redisClient.Do("SISMEMBER", cSet.itemLiked, userId); err != nil {
	//	return err
	//}
	//if rs == 0 {
	//	redisClient.Do("ZINCRBY", cSet.mostLiked, 1, itemId)
	//}

	if _, err := redisClient.Do("SADD", cSet.userLiked, itemId); err != nil {
		return err
	}

	if _, err := redisClient.Do("SADD", cSet.itemLiked, userId); err != nil {
		return err
	}

	return nil
}

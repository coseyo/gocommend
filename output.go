package gocommend

import "github.com/garyburd/redigo/redis"

type Output struct {
	collection string
	userId     string
	itemId     string
	recNum     int
	cSet       collectionSet
}

func (this *Output) Init(collection string, userId string, itemId string, recNum int) error {
	this.collection = collection
	this.userId = userId
	this.itemId = itemId
	this.recNum = recNum
	this.cSet = collectionSet{}
	this.cSet.init(collection)
	return nil
}

func (this *Output) toStrings(arrayInterface []interface{}) (strings []string) {
	for _, rs := range arrayInterface {
		s, _ := redis.String(rs, nil)
		strings = append(strings, s)
	}
	return
}

func (this *Output) RecommendedItem() ([]string, error) {
	arrayInterface, err := redis.Values(redisClient.Do("ZREVRANGE", this.cSet.recommendedItem(this.userId), 0, this.recNum))
	if err != nil {
		return nil, err
	}
	return this.toStrings(arrayInterface), err
}

//func (this *Output) MostLiked() ([]string, error) {
//	cSet := collectionSet{}
//	cSet.init(this.Collection)
//	arrayInterface, err := redis.Values(redisClient.Do("ZREVRANGE", cSet.recommendedItem(this.UserId), 0, this.RecNum))
//	if err != nil {
//		return nil, err
//	}
//	return this.toStrings(arrayInterface), err
//}

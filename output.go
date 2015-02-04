package gocommend

import "github.com/garyburd/redigo/redis"

type Output struct {
	Collection string
	UserId     string
	ItemId     string
	RecNum     int
}

func (this *Output) toStrings(arrayInterface []interface{}) (strings []string) {
	for _, rs := range arrayInterface {
		s, _ := redis.String(rs, nil)
		strings = append(strings, s)
	}
	return
}

func (this *Output) RecommendedItem() ([]string, error) {
	cSet := collectionSet{}
	cSet.init(this.Collection)
	arrayInterface, err := redis.Values(redisClient.Do("ZREVRANGE", cSet.recommendedItem(this.UserId), 0, this.RecNum))
	if err != nil {
		return nil, err
	}
	return this.toStrings(arrayInterface), err
}

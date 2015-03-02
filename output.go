package gocommend

import "github.com/garyburd/redigo/redis"

// output, get data what you want
type Output struct {
	recNum int
	cSet   collectionSet
}

// init the params and set the cSet
func (this *Output) Init(collection string, recNum int) error {
	if collection == "" {
		return gocommendError{emptyCollection}
	}
	this.recNum = recNum
	this.cSet = collectionSet{}
	this.cSet.init(collection)
	return nil
}

// convert interface slice to string slice
func (this *Output) toStrings(arrayInterface []interface{}) (strings []string) {
	for _, rs := range arrayInterface {
		s, _ := redis.String(rs, nil)
		strings = append(strings, s)
	}
	return
}

// get recommend items for user
func (this *Output) SimilarItemForUser(userId string) ([]string, error) {
	arrayInterface, err := redis.Values(redisClient.Do("ZREVRANGE", this.cSet.recommendedItem(userId), 0, this.recNum))
	if err != nil {
		return nil, err
	}
	return this.toStrings(arrayInterface), err
}

// get recommend items by item similarty
func (this *Output) SimilarItemForItem(itemId string) ([]string, error) {
	arrayInterface, err := redis.Values(redisClient.Do("ZREVRANGE", this.cSet.itemSimilarity(itemId), 0, this.recNum))
	if err != nil {
		return nil, err
	}
	return this.toStrings(arrayInterface), err
}

// get the best rated items
func (this *Output) BestRated() ([]string, error) {
	arrayInterface, err := redis.Values(redisClient.Do("ZREVRANGE", this.cSet.scoreRank, 0, this.recNum))
	if err != nil {
		return nil, err
	}
	return this.toStrings(arrayInterface), err
}

func (this *Output) MostLiked() ([]string, error) {
	arrayInterface, err := redis.Values(redisClient.Do("ZREVRANGE", this.cSet.mostLiked, 0, this.recNum))
	if err != nil {
		return nil, err
	}
	return this.toStrings(arrayInterface), err
}

func (this *Output) MostSimilarUsers(userId string) ([]string, error) {
	arrayInterface, err := redis.Values(redisClient.Do("ZREVRANGE", this.cSet.userSimilarity(userId), 0, this.recNum))
	if err != nil {
		return nil, err
	}
	return this.toStrings(arrayInterface), err
}

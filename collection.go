package gocommend

type collectionSet struct {
	collectionPrefix string
	scoreBoard       string
	mostLiked        string
	mostDisliked     string
}

func (c *collectionSet) init(collection string) {
	c.collectionPrefix = collection
	c.scoreBoard = c.collectionPrefix + ":scoreBoard"
	c.mostLiked = c.collectionPrefix + ":mostLiked"
	c.mostDisliked = c.collectionPrefix + ":mostDisliked"
}

func (c *collectionSet) userLiked(userId string) string {
	return c.collectionPrefix + ":" + userId + ":" + "userLiked"
}

func (c *collectionSet) itemLiked(itemId string) string {
	return c.collectionPrefix + ":" + itemId + ":" + "itemLiked"
}

func (c *collectionSet) userDisliked(userId string) string {
	return c.collectionPrefix + ":" + userId + ":" + "userDisliked"
}

func (c *collectionSet) itemDisliked(itemId string) string {
	return c.collectionPrefix + ":" + itemId + ":" + "itemDisliked"
}

func (c *collectionSet) userSimilarity(userId string) string {
	return c.collectionPrefix + ":" + userId + ":" + "userSimilarity"
}

func (c *collectionSet) itemSimilarity(itemId string) string {
	return c.collectionPrefix + ":" + itemId + ":" + "itemSimilarity"
}

func (c *collectionSet) userTemp(userId string) string {
	return c.collectionPrefix + ":" + userId + ":" + "userTemp"
}

func (c *collectionSet) itemTemp(itemId string) string {
	return c.collectionPrefix + ":" + itemId + ":" + "itemTemp"
}

func (c *collectionSet) userTempDiff(userId string) string {
	return c.collectionPrefix + ":" + userId + ":" + "userTempDiff"
}

func (c *collectionSet) itemTempDiff(userId string) string {
	return c.collectionPrefix + ":" + userId + ":" + "itemTempDiff"
}

func (c *collectionSet) recommendedItem(userId string) string {
	return c.collectionPrefix + ":" + userId + ":" + "recommendedItem"
}

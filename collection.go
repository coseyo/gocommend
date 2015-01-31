package gocommend

type collectionSet struct {
	mostLiked       string
	mostDisliked    string
	userLiked       string
	itemLiked       string
	userDisliked    string
	itemDisliked    string
	userSimilarity  string
	itemSimilarity  string
	userTemp        string
	itemTemp        string
	userTempDiff    string
	itemTempDiff    string
	recommendedItem string
}

func initSet(i *Input) (*collectionSet, error) {
	c := new(collectionSet)
	c.mostLiked = i.Collection + ":mostLiked"
	c.mostDisliked = i.Collection + ":mostDisliked"
	c.userLiked = i.Collection + ":" + i.UserId + ":" + "userLiked"
	c.itemLiked = i.Collection + ":" + i.ItemId + ":" + "itemLiked"
	c.userDisliked = i.Collection + ":" + i.UserId + ":" + "userDisliked"
	c.itemDisliked = i.Collection + ":" + i.ItemId + ":" + "itemDisliked"
	c.userSimilarity = i.Collection + ":" + i.UserId + ":" + "userSimilarity"
	c.itemSimilarity = i.Collection + ":" + i.ItemId + ":" + "itemSimilarity"
	c.userTemp = i.Collection + ":" + i.UserId + ":" + "userTemp"
	c.itemTemp = i.Collection + ":" + i.ItemId + ":" + "itemTemp"
	c.userTempDiff = i.Collection + ":" + i.UserId + ":" + "userTempDiff"
	c.itemTempDiff = i.Collection + ":" + i.UserId + ":" + "itemTempDiff"
	c.recommendedItem = i.Collection + ":" + i.UserId + ":" + "recommendedItem"

	return c, nil
}

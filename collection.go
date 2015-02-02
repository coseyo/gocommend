package gocommend

type collectionSet struct {
	collectionPrefix string
	scoreRank        string
	mostLiked        string
	mostDisliked     string
}

func (this *collectionSet) init(collection string) {
	this.collectionPrefix = collection
	this.scoreRank = this.collectionPrefix + ":scoreRank"
	this.mostLiked = this.collectionPrefix + ":mostLiked"
	this.mostDisliked = this.collectionPrefix + ":mostDisliked"
}

func (this *collectionSet) userLiked(userId string) string {
	return this.collectionPrefix + ":" + userId + ":" + "userLiked"
}

func (this *collectionSet) itemLiked(itemId string) string {
	return this.collectionPrefix + ":" + itemId + ":" + "itemLiked"
}

func (this *collectionSet) userDisliked(userId string) string {
	return this.collectionPrefix + ":" + userId + ":" + "userDisliked"
}

func (this *collectionSet) itemDisliked(itemId string) string {
	return this.collectionPrefix + ":" + itemId + ":" + "itemDisliked"
}

func (this *collectionSet) userSimilarity(userId string) string {
	return this.collectionPrefix + ":" + userId + ":" + "userSimilarity"
}

func (this *collectionSet) itemSimilarity(itemId string) string {
	return this.collectionPrefix + ":" + itemId + ":" + "itemSimilarity"
}

func (this *collectionSet) userTemp(userId string) string {
	return this.collectionPrefix + ":" + userId + ":" + "userTemp"
}

func (this *collectionSet) itemTemp(itemId string) string {
	return this.collectionPrefix + ":" + itemId + ":" + "itemTemp"
}

func (this *collectionSet) userTempDiff(userId string) string {
	return this.collectionPrefix + ":" + userId + ":" + "userTempDiff"
}

func (this *collectionSet) itemTempDiff(userId string) string {
	return this.collectionPrefix + ":" + userId + ":" + "itemTempDiff"
}

func (this *collectionSet) recommendedItem(userId string) string {
	return this.collectionPrefix + ":" + userId + ":" + "recommendedItem"
}

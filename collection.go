package gocommend

// redis key set
type collectionSet struct {
	collectionPrefix string
	scoreRank        string
	mostLiked        string
	mostDisliked     string
	allItem          string
	allUser          string
}

func (this *collectionSet) init(collection string) {
	this.collectionPrefix = DB_PREFIX + ":" + collection
	this.scoreRank = this.collectionPrefix + ":scoreRank"
	this.mostLiked = this.collectionPrefix + ":mostLiked"
	this.mostDisliked = this.collectionPrefix + ":mostDisliked"
	this.allItem = this.collectionPrefix + ":allItem"
	this.allUser = this.collectionPrefix + ":allUser"
}

func (this *collectionSet) userLiked(userId string) string {
	return this.collectionPrefix + ":" + "userLiked" + ":" + userId
}

func (this *collectionSet) itemLiked(itemId string) string {
	return this.collectionPrefix + ":" + "itemLiked" + ":" + itemId
}

func (this *collectionSet) userDisliked(userId string) string {
	return this.collectionPrefix + ":" + "userDisliked" + ":" + userId
}

func (this *collectionSet) itemDisliked(itemId string) string {
	return this.collectionPrefix + ":" + "itemDisliked" + ":" + itemId
}

func (this *collectionSet) userSimilarity(userId string) string {
	return this.collectionPrefix + ":" + "userSimilarity" + ":" + userId
}

func (this *collectionSet) itemSimilarity(itemId string) string {
	return this.collectionPrefix + ":" + "itemSimilarity" + ":" + itemId
}

func (this *collectionSet) userTemp(userId string) string {
	return this.collectionPrefix + ":" + "userTemp" + ":" + userId
}

func (this *collectionSet) itemTemp(itemId string) string {
	return this.collectionPrefix + ":" + "itemTemp" + ":" + itemId
}

func (this *collectionSet) userTempDiff(userId string) string {
	return this.collectionPrefix + ":" + "userTempDiff" + ":" + userId
}

func (this *collectionSet) itemTempDiff(userId string) string {
	return this.collectionPrefix + ":" + "itemTempDiff" + ":" + userId
}

func (this *collectionSet) recommendedItem(userId string) string {
	return this.collectionPrefix + ":" + "recommendedItem" + ":" + userId
}

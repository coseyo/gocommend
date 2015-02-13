package gocommend

const (
	// the num of most similar users to get for calculation the recommend item
	MAX_NEIGHBORS = 10

	// max recommend item num to store
	MAX_RECOMMEND_ITEM = 30

	DB_PREFIX = "gocommend"

	MAX_SIMILARITY_ITEM = 100

	MAX_SIMILARITY_USER = 100

	localRedisURL = "192.168.1.7"

	localRedisPort = "6379"

	remoteRedisURL = "10.20.187.251"

	remoteRedisPort = "11311"

	localStartup = false
)

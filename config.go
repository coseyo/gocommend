package gocommend

const (
	// the num of most similar users to get for calculation the recommend item
	MAX_NEIGHBORS = 10

	// max recommend item num to store
	MAX_RECOMMEND_ITEM = 30

	// redis key prefix
	DB_PREFIX = "gocommend"

	MAX_SIMILARITY_ITEM = 100

	MAX_SIMILARITY_USER = 100

	LOCAL_REDIS_HOST = "192.168.1.7"

	LOCAL_REDIS_PORT = "6379"

	REMOTE_REDIS_HOST = "10.20.187.251"

	REMOTE_REDIS_PORT = "11311"

	LOCAL_STARTUP = false
)

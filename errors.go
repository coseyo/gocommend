package gocommend

const (
	emptyCollection = "Empty collection."
	missingHeader   = "Missing required header: %q"
)

type gocommendError struct {
	Message string
}

func (e gocommendError) Error() string {
	return e.Message
}

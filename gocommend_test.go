package gocommend

import (
	"reflect"
	"testing"

	"github.com/garyburd/redigo/redis"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func Test_redisTest(t *testing.T) {
	redisClient.Do("SET", "aaa", 123)
	a, err := redis.Int(redisClient.Do("GET", "aaa"))
	expect(t, err, nil)
	expect(t, a, 123)
}

func Test_importPoll(t *testing.T) {
	collection := "rec_test2"
	i := Input{}
	i.Init(collection)
	err := i.ImportPoll("u1", "i1")
	expect(t, err, nil)
	i.ImportPoll("u1", "i2")
	i.ImportPoll("u1", "i3")
	i.ImportPoll("u2", "i1")
	i.ImportPoll("u2", "i2")
}

func Test_updatePoll(t *testing.T) {
	collection := "rec_test2"
	i := Input{}
	i.Init(collection)
	err := i.UpdatePoll("u1", "")
	expect(t, err, nil)
	i.UpdatePoll("u2", "")
}

func Test_updateAllPoll(t *testing.T) {
	collection := "rec_test2"
	i := Input{}
	i.Init(collection)
	err := i.UpdateAllPoll()
	expect(t, err, nil)
}

func Test_RecommendItem(t *testing.T) {
	collection := "rec_test2"
	recNum := 10
	o := Output{}
	o.Init(collection, recNum)
	items, err := o.RecommendItemForUser("u2")
	expect(t, err, nil)
	expect(t, items[0], "i3")
}

gocommend
===========
recommend system for golang.

Developing...

## Example code below

```go

package main

import (
	"fmt"
	"gocommend"
	"log"
	"os"
	"strconv"
)

func main() {

	//gocommend.Redistest()

	argNum := len(os.Args)

	handle := os.Args[1]
	collection := os.Args[2]

	switch handle {
	case "import":
		if argNum != 6 {
			fmt.Println("num of input params shuold be 5")
			return
		}
		userId := os.Args[3]
		itemId := os.Args[4]
		rate, _ := strconv.Atoi(os.Args[5])
		i := gocommend.Input{
			Collection: collection,
			UserId:     userId,
			ItemId:     itemId,
			Rate:       rate,
		}
		i.ImportPoll()

	case "update":
		userId := os.Args[3]
		itemId := os.Args[4]
		i := gocommend.Input{
			Collection: collection,
			UserId:     userId,
			ItemId:     itemId,
		}
		err := i.UpdatePoll()
		if err != nil {
			log.Println(err)
		}

	case "rec":
		userId := os.Args[3]
		o := gocommend.Output{
			Collection: collection,
			UserId:     userId,
			RecNum:     10,
		}
		rs, err := o.RecommendedItem()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(rs)
	}
}


````
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
		gocommend.ImportPoll(&gocommend.Input{
			Collection: collection,
			UserId:     userId,
			ItemId:     itemId,
			Rate:       rate,
		})

	case "update":
		userId := os.Args[3]
		itemId := os.Args[4]
		err := gocommend.UpdatePoll(&gocommend.Input{
			Collection: collection,
			UserId:     userId,
			ItemId:     itemId,
		})
		if err != nil {
			log.Println(err)
		}
	}
}



````
package main

import (
	"fmt"
	"gocommend"
	"log"
	"os"
)

func main() {

	//gocommend.Redistest()

	argNum := len(os.Args)

	handle := os.Args[1]
	collection := os.Args[2]

	switch handle {
	case "importPoll":
		if argNum != 5 {
			fmt.Println("num of input params shuold be 5")
			return
		}
		userId := os.Args[3]
		itemId := os.Args[4]
		//rate, _ := strconv.Atoi(os.Args[5])
		i := gocommend.Input{}
		i.Init(collection)
		i.ImportPoll(userId, itemId)

	case "updatePoll":
		userId := os.Args[3]
		//itemId := os.Args[4]
		i := gocommend.Input{}
		i.Init(collection)
		err := i.UpdatePoll(userId, "")
		if err != nil {
			log.Println(err)
		}

	case "recommendForUser":
		userId := os.Args[3]
		//itemId := os.Args[4]
		recNum := 10
		o := gocommend.Output{}
		o.Init(collection, recNum)
		rs, err := o.RecommendItemForUser(userId)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(rs)

	case "recommendForItem":
		itemId := os.Args[3]
		recNum := 10
		o := gocommend.Output{}
		o.Init(collection, recNum)
		rs, err := o.RecommendItemForItem(itemId)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(rs)
	}
}

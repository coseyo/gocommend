gocommend
===========
recommend system for golang.

Developing...

## CLI example code below

```go

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
		rs, err := o.SimilarItemForUser(userId)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(rs)
	}
	
	case "recommendForItem":
		itemId := os.Args[3]
		recNum := 10
		o := gocommend.Output{}
		o.Init(collection, recNum)
		rs, err := o.SimilarItemForItem(itemId)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(rs)
}



````

## HTTP example code below

```go

package main

import (
	"encoding/json"
	"gocommend"
	"log"
	"net/http"
	"strconv"
)

type commendServer struct {
	w        http.ResponseWriter
	req      *http.Request
	postData map[string][]string
}

func (this *commendServer) init(w http.ResponseWriter, req *http.Request) (err string) {
	this.w = w
	this.req = req
	this.req.ParseForm()
	this.postData = this.req.PostForm
	if len(this.postData) == 0 {
		err = "No post data"
	}
	return
}

func (this *commendServer) responseJson(result string, data interface{}, msg string) {
	this.w.Header().Set("content-type", "application/json")
	jsonData := map[string]interface{}{
		"result": result,
		"data":   data,
		"msg":    msg,
	}
	rs, _ := json.Marshal(jsonData)
	this.w.Write(rs)
}

func (this *commendServer) getParam(key string, allowNull bool) (value string, err string) {
	valueArray, exist := this.postData[key]
	if allowNull == true {
		if exist == false {
			return "", ""
		}
		err = ""
	} else {
		if exist == false {
			err = " No key " + key
			return
		}
		if valueArray[0] == "" {
			err = " empty value " + key
		}
	}
	value = valueArray[0]
	return
}

func importPollHandler(w http.ResponseWriter, req *http.Request) {
	s := commendServer{}
	if err := s.init(w, req); err != "" {
		s.responseJson("error", "", err)
		return
	}

	collection, err1 := s.getParam("collection", false)
	userId, err2 := s.getParam("userId", false)
	itemId, err3 := s.getParam("itemId", false)
	if err1 != "" || err2 != "" || err3 != "" {
		s.responseJson("error", "", err1+err2+err3)
		return
	}

	i := gocommend.Input{}
	i.Init(collection)
	if err := i.ImportPoll(userId, itemId); err != nil {
		s.responseJson("error", "", err.Error())
		return
	}
	s.responseJson("ok", "", "")
}

func updatePollHandler(w http.ResponseWriter, req *http.Request) {
	s := commendServer{}
	if err := s.init(w, req); err != "" {
		s.responseJson("error", "", err)
		return
	}

	collection, err1 := s.getParam("collection", false)
	userId, err2 := s.getParam("userId", false)
	itemId, err3 := s.getParam("itemId", true)
	if err1 != "" || err2 != "" || err3 != "" {
		s.responseJson("error", "", err1+err2+err3)
		return
	}

	i := gocommend.Input{}
	i.Init(collection)
	if err := i.UpdatePoll(userId, itemId); err != nil {
		s.responseJson("error", "", err.Error())
		return
	}
	s.responseJson("ok", "", "")
}

func similarItemForUserHandler(w http.ResponseWriter, req *http.Request) {
	s := commendServer{}
	if err := s.init(w, req); err != "" {
		s.responseJson("error", "", err)
		return
	}

	collection, err1 := s.getParam("collection", false)
	userId, err2 := s.getParam("userId", false)
	num, err3 := s.getParam("num", true)
	if err1 != "" || err2 != "" || err3 != "" {
		s.responseJson("error", "", err1+err2+err3)
		return
	}

	recNum := 10
	if num != "" {
		recNum, _ = strconv.Atoi(num)
	}
	o := gocommend.Output{}
	o.Init(collection, recNum)
	rs, err := o.SimilarItemForUser(userId)
	log.Println(rs)
	if err != nil {
		s.responseJson("error", "", err.Error())
		return
	}
	s.responseJson("ok", rs, "")
}

func similarItemForItemHandler(w http.ResponseWriter, req *http.Request) {
	s := commendServer{}
	if err := s.init(w, req); err != "" {
		s.responseJson("error", "", err)
		return
	}

	collection, err1 := s.getParam("collection", false)
	itemId, err2 := s.getParam("itemId", false)
	num, err3 := s.getParam("num", true)
	if err1 != "" || err2 != "" || err3 != "" {
		s.responseJson("error", "", err1+err2+err3)
		return
	}

	recNum := 10
	if num != "" {
		recNum, _ = strconv.Atoi(num)
	}

	o := gocommend.Output{}
	o.Init(collection, recNum)
	rs, err := o.SimilarItemForItem(itemId)
	if err != nil {
		s.responseJson("error", "", err.Error())
		return
	}
	s.responseJson("ok", rs, "")
}

func main() {

	http.HandleFunc("/importPoll", importPollHandler)
	http.HandleFunc("/updatePoll", updatePollHandler)
	http.HandleFunc("/similarItemForUser", similarItemForUserHandler)
	http.HandleFunc("/similarItemForItem", similarItemForItemHandler)

	http.ListenAndServe(":8888", nil)
}


````
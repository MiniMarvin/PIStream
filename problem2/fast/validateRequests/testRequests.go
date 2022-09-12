package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const wordLen = 21
const requestLen = 50

type PiRequest struct {
	Content string
}

func getJson(url string) (error, string) {
	r, err := http.Get(url)
	if err != nil {
		return err, ""
	}

	if r.Body != nil {
		defer r.Body.Close()
	}

	body, readErr := ioutil.ReadAll(r.Body)

	if readErr != nil {
		return readErr, ""
	}

	// 	fmt.Println("body: ", body)

	target := PiRequest{Content: ""}
	jsonErr := json.Unmarshal(body, &target)
	if jsonErr != nil {
		return jsonErr, ""
	}

	return nil, target.Content
}

func query(begin int, count int, ch chan<- bool) {
	path := fmt.Sprintf("https://api.pi.delivery/v1/pi?start=%d&numberOfDigits=%d", begin, count)
	// 	fmt.Println(path)
	err, target := getJson(path)
	// 	fmt.Println("target: ", target)
	if err != nil {
		ch <- false
		return
	}
	ch <- len(target) == requestLen
}

func main() {
	begin := 0
	step := requestLen
	requestLimit := 100000
	success := 0
	total := 0

	c := make(chan bool)

	for i := 0; i < requestLimit; i++ {
		go query(begin, step, c)
		begin += step - wordLen
	}

	for i := 0; i < requestLimit; i++ {
		v := <-c
		if v {
			success += 1
		}
		total += 1
	}

	fmt.Println("success: ", success, "/", total)
}

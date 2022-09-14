package main

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lukechampine.com/uint128"
	"net/http"
	"os"
	"primality"
	"strconv"
)

const HIGH_LIMIT uint64 = 200000000000 // 100Bi
const START_INDEX uint64 = 100000000000
const STEP uint64 = 1000
const SEQUENCE_LEN uint64 = 21
const POOL_LIMIT uint64 = 10000

var finished chan string
var resultChan chan Item
var finishLock bool

func resultManager() {
	resultHeap := make(PriorityQueue, 0)
	for {
		item := <-resultChan
		finishLock = true

		pushItem := &Item{
			value:    item.value,
			priority: item.priority,
		}
		heap.Push(&resultHeap, pushItem)
		fmt.Println("found new result!!! - " + item.value)
		minVal := resultHeap.Top()

		success := false
		for !success {
			file, err := os.OpenFile("foundNumbers.txt",
				os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			file1, err1 := os.OpenFile("detectedNumber.txt",
				os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil || err1 != nil {
				continue
			}

			_, err2 := file.WriteString(fmt.Sprintf("[%d] %s\n", item.priority, item.value))
			_, err3 := file1.WriteString(minVal)
			if err2 == nil && err3 == nil {
				success = true
			}

			file.Close()
			file1.Close()
		}
		finishLock = false
	}
}

func appendToResult(s string, idx uint64) {
	// TODO: check if file exists
	resultChan <- Item{
		value:    s,
		priority: idx,
	}
}

func writeLog(
	startChan <-chan uint64,
	endChan <-chan uint64) {

	fmt.Println("[writeLog] Init heap structure")
	// TODO: write into afile the last visited index
	receivedHeap := &IntHeap{}
	finishedHeap := &IntHeap{}

	var minVal uint64
	var maxVal uint64
	minVal = 0
	maxVal = 0

	fmt.Println("[writeLog] Long poll on channels...")
	for {
		select {
		case rcvd := <-startChan:
			{
				heap.Push(receivedHeap, rcvd)
				if rcvd >= maxVal {
					minVal, _ = receivedHeap.Top()
					maxVal = rcvd
					// 	fmt.Println("[writeLog][rcv] searching in range: [",
					// 		minVal, ", ", maxVal+STEP, "]")
				}
			}
		case finished := <-endChan:
			{
				heap.Push(finishedHeap, finished)
				t1, hast1 := (*receivedHeap).Top()
				t2, hast2 := (*finishedHeap).Top()
				if hast1 && hast2 && t1 == t2 {
					receivedHeap.Pop()
					finishedHeap.Pop()
					s := strconv.FormatInt(int64(minVal), 10)
					f, _ := os.Create("lastFinished.txt")
					f.WriteString(s + "\n")
					f.Close()

					minVal, _ = receivedHeap.Top()
					// 	fmt.Println("[writeLog][end] searching in range: [",
					// 		minVal, ", ", maxVal+STEP, "]")
				}
			}
		}
	}

	fmt.Println("[writeLog] Finished Long poll")
}

func checkPalindrome(s string) bool {
	isPalindrome := true
	limit := len(s) / 2
	for i := 0; i < limit; i++ {
		if s[i] != s[len(s)-1-i] {
			isPalindrome = false
			break
		}
	}
	return isPalindrome
}

func searchInIdx(
	startIdx uint64,
	count uint64,
	successChan chan<- uint64,
	failureChan chan<- uint64) {
	path := fmt.Sprintf(
		"https://api.pi.delivery/v1/pi?start=%d&numberOfDigits=%d",
		startIdx, count)
	err, digits := getJson(path)
	if err != nil {
		failureChan <- startIdx
		return
	}

	var i uint64
	for i = 0; i < count-SEQUENCE_LEN+1; i++ {
		s := digits[i : i+SEQUENCE_LEN]
		if checkPalindrome(s) {
			fmt.Println(fmt.Sprintf("[%d][%d] palindrome try", startIdx+i, s))
			p, err := uint128.FromString(s)
			if err != nil {
				failureChan <- startIdx
				return
			}

			if primality.PrimeCheck(p, false) {
				appendToResult(s, startIdx+i)
				return
			}
		}
	}

	successChan <- startIdx
}

func nextStep(last uint64) uint64 {
	return last + STEP - SEQUENCE_LEN + 1
}

func startDispatching(
	startIdx uint64,
	successChan chan uint64,
	failureChan chan uint64) uint64 {
	var last uint64
	last = startIdx
	var i uint64
	for i = 0; i < POOL_LIMIT; i++ {
		go searchInIdx(last, STEP, successChan, failureChan)
		last = nextStep(last)
	}
	return last
}

func keepDispatching(
	startIdx uint64,
	successChan chan uint64,
	failureChan chan uint64,
	logStartChan chan<- uint64,
	logEndChan chan<- uint64) {
	var last uint64
	last = startIdx

	fmt.Println("[keepDispatching] init loop over channels")
	for last < HIGH_LIMIT {
		select {
		case finishedIdx := <-successChan:
			{
				// fmt.Println(fmt.Sprintf("[keepDispatching][%d] received finished idx", finishedIdx))
				logEndChan <- finishedIdx
				last = nextStep(last)
				logStartChan <- last
				fmt.Println(fmt.Sprintf("[keepDispatching][%d][%d] starting next query", finishedIdx, last))
				go searchInIdx(last, STEP, successChan, failureChan)
			}
		case failedIdx := <-failureChan:
			{
				fmt.Println(fmt.Sprintf("[keepDispatching][%d] received failure idx", failedIdx))
				go searchInIdx(failedIdx, STEP, successChan, failureChan)
			}
		}
	}
	fmt.Println("[keepDispatching] reached limit value, quiting...")

	finished <- ""
}

// Handles dispatching of starting indexes
// and index retries
func dispatcher(
	startIdx uint64) {

	fmt.Println("[dispatcher] creating status channels...")
	successChan := make(chan uint64, POOL_LIMIT)
	failureChan := make(chan uint64, POOL_LIMIT)

	fmt.Println("[dispatcher] starting dispatches...")
	last := startDispatching(startIdx, successChan, failureChan)

	fmt.Println("[dispatcher] creating log channels...")
	logStartChan := make(chan uint64, POOL_LIMIT)
	logEndChan := make(chan uint64, POOL_LIMIT)

	fmt.Println("[dispatcher] starting log thread...")
	go writeLog(logStartChan, logEndChan)

	fmt.Println("[dispatcher] entering dispatch loop...")
	keepDispatching(
		last,
		successChan, failureChan,
		logStartChan, logEndChan)
	fmt.Println("[dispatcher] leaving dispatch loop...")
}

type PiRequest struct {
	Content string
}

func getJson(url string) (error, string) {
	r, err := http.Get(url)
	if err != nil {
		return err, ""
	}

	body, readErr := ioutil.ReadAll(r.Body)
	r.Body.Close()

	if readErr != nil {
		return readErr, ""
	}

	target := PiRequest{Content: ""}
	jsonErr := json.Unmarshal(body, &target)
	if jsonErr != nil {
		return jsonErr, ""
	}

	return nil, target.Content
}

func main() {
	// TODO: use the file from log to reload the start
	// index from there
	finished = make(chan string)
	resultChan = make(chan Item)
	go resultManager()
	go dispatcher(START_INDEX)
	<-finished

	fmt.Println("Results finished, waiting lock release...")
	for finishLock {
	}
	fmt.Println("Program finished check possible answers on file 'detectedNumber.txt'")
}

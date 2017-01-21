package main

import (
	"fmt"
	"sync"
	"time"
)

var queue = make(chan int, 10)
var wg *sync.WaitGroup = new(sync.WaitGroup)

func consumer(q chan int) {
	wg.Add(1)
	defer wg.Done()
	//	for task := range q {
	//		time.Sleep(time.Second)
	//		fmt.Println(task)
	//	}
Loop:
	for true {
		select {
		case <-time.After(time.Second * 2):
			break Loop
		case task := <-q:
			fmt.Println(task)
		}
	}

}

func delaySend(item int, q chan int) {
	fmt.Println("delay send before")
	q <- item
	fmt.Println("delay send done:", item)
}

func main() {
	go consumer(queue)
	for i := 0; i < 30; i++ {
		//time.Sleep(time.Second)
		select {
		case queue <- i:
			fmt.Println("main:put directly", i)
		default:
			go delaySend(i, queue)
		}

		//fmt.Println("main:", i)

	}
	//close(queue)
	wg.Wait()
	fmt.Println("Hello World")
}

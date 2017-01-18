package main

import (
	"fmt"
	"sync"
	"time"
    _ "net/http/pprof"
	"net/http"
)

var orderNum = 0

var locker  = &sync.Mutex{}

func getNum() int {
	locker.Lock()
	defer locker.Unlock()
	orderNum += 1
	return orderNum
}

func panicTest(){
	defer func() {
		if err:= recover(); err != nil{
			fmt.Println(err)
			go panicTest()
		}
	}()
    num := getNum()
	for i:=0; ;i++ {
		fmt.Println("my num is ", num, time.Now().Unix())
		time.Sleep(time.Second)
		if i > 2 && num < 5 {
				panic(fmt.Sprintln(num, "trigger the panic"))
		}
	}
}

func main() {
	fmt.Println("hello")
	go http.ListenAndServe(":8080", http.DefaultServeMux)
	go panicTest()
	time.Sleep(time.Second * 200)
}

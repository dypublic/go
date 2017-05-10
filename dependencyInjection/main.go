package main

import (
	"fmt"

	"./depandency"
)

func main() {
	do("first")
}

func do(name string) {
	depandency.DepFirst(name)
	fmt.Println("finish")
}

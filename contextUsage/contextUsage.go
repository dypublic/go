package main

import (
	"context"
	"fmt"
	"time"
	"log"
	_ "net/http/pprof"
	"net/http"
)

func heavyWork(count int){
	total := 0
	for i:=0; i < count; i++{
		total += i
	}
}

func processMessage(ctx context.Context, message string) error{
	for char := range message{
		fmt.Println(char)
		heavyWork(strconv.)
		select {
		case <-ctx.Done():
			fmt.Println("some one cancel it")
			return ctx.Err()
		default:
			fmt.Println("one pass")
		}

	}
	return nil
}

type Content struct {
	ctx context.Context
	message string
}

func process1(contentQ chan Content)  {
	for content := range contentQ{
		go processMessage(content.ctx, content.message)
	}
}

func main() {
	go func() {
        log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()
	contentQueue := make(chan Content)
	go process1(contentQueue)
	ctx, cancelFunc := context.WithCancel(context.Background())
	cont := Content{ctx, "test12345678test12345678test12345678test12345678"}
	contentQueue <- cont
	<-time.After(time.Second*3000)
	cancelFunc()
	<-time.After(time.Second*3)



}

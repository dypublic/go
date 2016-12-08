package main

import (
	"fmt"
	"sync"
	"time"
)

type task interface {
	run() bool
}

type pool struct {
	wg         *sync.WaitGroup
	size       int
	taskChan   chan task
	resultChan chan bool
}

type sleepTask struct {
	sleep  int
	taskId int
}

func (t *sleepTask) run() bool {
	fmt.Println("taskid:", t.taskId, "start!")
	time.Sleep(time.Duration(t.sleep)*time.Second)
	fmt.Println("taskid:", t.taskId, "end after:", t.sleep)
	return true
}

func engine(id int, taskChan chan task, resultChan chan bool, wg *sync.WaitGroup) {
    defer wg.Done()
	for {
		select {
		case t, ok := <-taskChan:
			if !ok {
				fmt.Println("engine:", id, " reveive close")
                return
			}
			res := t.run()
			resultChan <- res
		}
	}
}

func createPool(number int) pool {
	taskPool := pool{
		size:       number,
		wg:         new(sync.WaitGroup),
		taskChan:   make(chan task, 200),
		resultChan: make(chan bool, 200),
	}

	for i := 0; i < taskPool.size; i++ {
		taskPool.wg.Add(1)
		go engine(i, taskPool.taskChan, taskPool.resultChan, taskPool.wg)
	}
	return taskPool
}

func createTestTaskList(number int) []task {
	taskList := make([]task, 0, number)
	for i := 0; i < number; i++ {
		var temp task = &sleepTask{sleep: i, taskId: i}
		taskList = append(taskList, temp)
	}
	return taskList
}

func collectResult(resultChan chan bool) {
	for {
		select {
		case <-resultChan:
			fmt.Println("get one result")
		}
	}
}

func main() {
	pool := createPool(10)
	list := createTestTaskList(20)
	for _, t := range list {
		pool.taskChan <- t
	}
	go collectResult(pool.resultChan)
	time.Sleep(50*time.Second)
	close(pool.taskChan)
	pool.wg.Wait()
	close(pool.resultChan)

}

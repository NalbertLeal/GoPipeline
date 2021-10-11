package main

import (
	"fmt"
	"time"

	pipeline "github.com/NalbertLeal/GoPipeline"
)

func producer(in chan *pipeline.Data) {
	for i := 0; i < 3; i++ {
		in <- pipeline.NewData(i, pipeline.Processing)
	}
	in <- pipeline.NewData(nil, pipeline.Completed)
}

func mult(v interface{}) (interface{}, error) {
	return v.(int) * 2, nil
}

func add(v interface{}) (interface{}, error) {
	return v.(int) + 1, nil
}

func onData(v interface{}) {
	fmt.Println("Data: ", v.(int))
}

func onError(e error) {
	fmt.Println("Error: ", e.Error())
}

func onComplete() {
	fmt.Println("Completed")
}

func main() {
	ch := make(chan interface{}, 3)
	ch <- 1
	ch <- 6
	ch <- 4

	p := pipeline.FromChannel(ch)
	p.Pipe(mult).Pipe(add)

	go func() {
		time.Sleep(time.Second)
		p.Close()
	}()

	p.Subscribe(onData, onError, onComplete)
}

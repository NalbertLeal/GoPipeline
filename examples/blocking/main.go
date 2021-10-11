package main

import (
	"fmt"

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

func main() {
	p := pipeline.NewPipeline(producer)
	p.Pipe(mult).Pipe(add)
	values, _ := p.Block()
	fmt.Println(values)
}

package pipeline

const (
	Processing = iota
	PipelineError
	Completed
)

type Producer = func(chan *Data)
type Stage = func(value interface{}) (interface{}, error)
type PipelineStatus = uint8
type OnData = func(value interface{})
type OnError = func(err error)
type OnComplete = func()

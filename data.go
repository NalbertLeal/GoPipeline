package pipeline

type Data struct {
	Value  interface{}
	Status PipelineStatus
}

func NewData(value interface{}, status PipelineStatus) *Data {
	return &Data{
		value,
		status,
	}
}

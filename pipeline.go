package pipeline

type Pipeline struct {
	stages []Stage
	in     chan *Data
	out    chan *Data
}

func FromProducer(producer Producer) *Pipeline {
	pipelineInput := make(chan *Data)

	go producer(pipelineInput)

	return &Pipeline{
		[]Stage{},
		pipelineInput,
		nil,
	}
}

func FromChannel(in chan interface{}) *Pipeline {
	pipelineInput := make(chan *Data)
	go fromChannelHelper(in, pipelineInput)

	return &Pipeline{
		[]Stage{},
		pipelineInput,
		nil,
	}
}

func fromChannelHelper(in chan interface{}, pipelineInput chan *Data) {
	for value := range in {
		pipelineInput <- NewData(value, Processing)
	}
}

func Just(values ...interface{}) *Pipeline {
	pipelineInput := make(chan *Data)
	go justHelper(values, pipelineInput)

	return &Pipeline{
		[]Stage{},
		pipelineInput,
		nil,
	}
}

func justHelper(values []interface{}, pipelineInput chan *Data) {
	for _, value := range values {
		pipelineInput <- NewData(value, Processing)
	}
	pipelineInput <- NewData(nil, Completed)
}

func (p *Pipeline) Pipe(stage Stage) *Pipeline {
	p.stages = append(p.stages, stage)
	return p
}

func (p *Pipeline) Subscribe(onData OnData, onError OnError, onComplete OnComplete) {
	p.chainStages()

	for data := range p.out {
		if data.Status == Processing {
			onData(data.Value)
		} else if data.Status == PipelineError {
			err := data.Value.(error)
			onError(err)
			p.closeInOut()
			break
		} else if data.Status == Completed {
			onComplete()
			p.closeInOut()
			break
		}
	}
}

func (p *Pipeline) Block() ([]interface{}, error) {
	p.chainStages()

	values := []interface{}{}
	for data := range p.out {
		if data.Status == Processing {
			values = append(values, data.Value)
		} else if data.Status == PipelineError {
			err := data.Value.(error)
			p.closeInOut()
			return values, err
		} else if data.Status == Completed {
			p.closeInOut()
			return values, nil
		}
	}
	return values, nil
}

func (p *Pipeline) Close() {
	p.in <- NewData(nil, Completed)
}

func (p *Pipeline) chainStages() {
	lastOut := p.in
	for i := 0; i < len(p.stages); i++ {
		stage := p.stages[i]
		currentOut := make(chan *Data)
		go p.runStage(lastOut, currentOut, stage)
		lastOut = currentOut
	}
	p.out = lastOut
}

func (p *Pipeline) runStage(in chan *Data, out chan *Data, stage Stage) {
	for data := range in {
		if data.Status == Processing {
			result, err := stage(data.Value)
			if err != nil {
				out <- NewData(result, PipelineError)
			} else {
				out <- NewData(result, Processing)
			}
		} else if data.Status == PipelineError {
			out <- data
			close(in)
			break
		} else if data.Status == Completed {
			out <- data
			close(in)
			break
		}
	}
}

func (p *Pipeline) closeInOut() {
	close(p.out)
}

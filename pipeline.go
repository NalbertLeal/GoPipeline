package pipeline

type Pipeline struct {
	stages []Stage
	in     chan *Data
	out    chan *Data
}

func NewPipeline(producer Producer) *Pipeline {
	pipelineInput := make(chan *Data)

	go producer(pipelineInput)

	return &Pipeline{
		[]Stage{},
		pipelineInput,
		nil,
	}
}

// func FromIterator() (*Pipeline, error) {
// 	return nil, nil
// }

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
		} else if data.Status == Completed {
			out <- data
			close(in)
		}
	}
}

func (p *Pipeline) closeInOut() {
	close(p.in)
	close(p.out)
}

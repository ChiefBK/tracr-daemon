package channels

type CollectorOutputProcessorInput struct {
	Output interface{}
	Key string
}

type ReceiverOutputProcessorInput struct {
	Output interface{}
	Key string
}

var CollectorProcessorChan = make(chan CollectorOutputProcessorInput)
var ReceiverProcessorChan = make(chan ReceiverOutputProcessorInput)

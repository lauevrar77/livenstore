package services

import (
	"sync"

	"github.com/oklog/ulid/v2"
)

type LinkEvent struct {
	StreamName string
	EventID    ulid.ULID
}

type StreamWriterContext struct {
	InputChan  chan LinkEvent
	OutputChan chan error
}

type SafeStreamWriter struct {
	BasePath      string
	WriterFactory StreamWriterFactory
	Writers       map[string]StreamWriterContext
	mut           sync.Mutex
}

func (sw *SafeStreamWriter) LinkEvent(streamName string, eventID ulid.ULID) error {
	var writerContext StreamWriterContext
	ok := false
	sw.mut.Lock()
	if writerContext, ok = sw.Writers[streamName]; !ok {
		writerContext = StreamWriterContext{
			InputChan:  make(chan LinkEvent, 0),
			OutputChan: make(chan error, 0),
		}
		sw.Writers[streamName] = writerContext
		writer := sw.WriterFactory(sw.BasePath)
		go func() {
			for event := range sw.Writers[streamName].InputChan {
				err := writer.LinkEvent(event.StreamName, event.EventID)
				sw.Writers[streamName].OutputChan <- err
			}
		}()
	}
	sw.mut.Unlock()

	writerContext.InputChan <- LinkEvent{
		StreamName: streamName,
		EventID:    eventID,
	}

	return <-writerContext.OutputChan
}

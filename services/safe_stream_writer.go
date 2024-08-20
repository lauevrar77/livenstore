package services

import (
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
	"livenstore.evrard.online/persistance"
)

var MAX_STREAM_WRITERS = 1000

type LinkEvent struct {
	StreamName string
	EventID    ulid.ULID
}

type StreamWriterContext struct {
	InputChan          chan LinkEvent
	OutputChan         chan error
	LastEventTimestamp int64
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
			InputChan:          make(chan LinkEvent, 0),
			OutputChan:         make(chan error, 0),
			LastEventTimestamp: time.Now().UnixMicro(),
		}
		sw.Writers[streamName] = writerContext
		writer := sw.WriterFactory(sw.BasePath)
		go func(inputChan chan LinkEvent, writer persistance.StreamWriter, outputChan chan error) {
			for event := range inputChan {
				err := writer.LinkEvent(event.StreamName, event.EventID)
				outputChan <- err
			}
		}(writerContext.InputChan, writer, writerContext.OutputChan)

		// Cleaning up the oldest stream writer if we have too many
		if len(sw.Writers) > MAX_STREAM_WRITERS {
			oldestTimestamp := time.Now().UnixMicro()
			oldestStreamName := ""
			for streamName, writerContext := range sw.Writers {
				if writerContext.LastEventTimestamp < oldestTimestamp {
					oldestTimestamp = writerContext.LastEventTimestamp
					oldestStreamName = streamName
				}
			}
			close(sw.Writers[oldestStreamName].InputChan)
			delete(sw.Writers, oldestStreamName)
		}
	}
	sw.mut.Unlock()

	writerContext.InputChan <- LinkEvent{
		StreamName: streamName,
		EventID:    eventID,
	}

	return <-writerContext.OutputChan
}

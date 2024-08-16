package services

import (
	"fmt"

	"github.com/oklog/ulid/v2"
	"livenstore.evrard.online/domain"
	"livenstore.evrard.online/persistance"
)

type EventWriterFactory func(string) persistance.EventWriter
type EventReaderFactory func(string) persistance.EventReader
type StreamWriterFactory func(string) persistance.StreamWriter

type WriteResult struct {
	Offset int64
	Error  error
}

type EventStore struct {
	BasePath         string
	Writer           persistance.EventWriter
	WriterChan       chan domain.Event
	WriterResultChan chan WriteResult
	ReaderFactory    EventReaderFactory
}

func NewEventStore(basePath string, writerFactory EventWriterFactory, readerFactory EventReaderFactory, streamWriterFactory StreamWriterFactory) EventStore {
	writer := writerFactory(basePath)
	writerChan := make(chan domain.Event, 0)
	writerResultChan := make(chan WriteResult, 0)

	streamWriter := streamWriterFactory(basePath)

	go func() {
		for event := range writerChan {
			offset, err := writer.WriteEvent(event)
			if err != nil {
				writerResultChan <- WriteResult{Offset: offset, Error: err}
			}
			err = streamWriter.LinkEvent(fmt.Sprintf("by_event_type_%s", event.Type), event.ID)
			writerResultChan <- WriteResult{Offset: offset, Error: err}
		}
	}()

	return EventStore{
		BasePath:         basePath,
		Writer:           writer,
		WriterChan:       writerChan,
		WriterResultChan: writerResultChan,
		ReaderFactory:    readerFactory,
	}
}

func (es *EventStore) WriteEvent(event domain.Event) (int64, error) {
	es.WriterChan <- event
	result := <-es.WriterResultChan
	return result.Offset, result.Error
}

func (es *EventStore) ReadEvent(eventID ulid.ULID) (*domain.Event, error) {
	reader := es.ReaderFactory(es.BasePath)
	return reader.ReadEvent(eventID)
}

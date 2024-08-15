package services

import (
	"github.com/oklog/ulid/v2"
	"livenstore.evrard.online/domain"
	"livenstore.evrard.online/persistance"
)

type EventWriterFactory func(string) persistance.EventWriter
type EventReaderFactory func(string) persistance.EventReader

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

func NewEventStore(basePath string, writerFactory EventWriterFactory, readerFactory EventReaderFactory) EventStore {
	writer := writerFactory(basePath)
	writerChan := make(chan domain.Event, 0)
	writerResultChan := make(chan WriteResult, 0)

	go func() {
		for event := range writerChan {
			offset, err := writer.WriteEvent(event)
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

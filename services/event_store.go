package services

import (
	"livenstore.evrard.online/domain"
	"livenstore.evrard.online/persistance"
)

type WriteResult struct {
	Offset int64
	Error  error
}

type EventStore struct {
	Writer           persistance.EventWriter
	WriterChan       chan domain.Event
	WriterResultChan chan WriteResult
}

func NewEventStore(writer persistance.EventWriter) EventStore {
	writerChan := make(chan domain.Event, 0)
	writerResultChan := make(chan WriteResult, 0)

	go func() {
		for event := range writerChan {
			offset, err := writer.WriteEvent(event)
			writerResultChan <- WriteResult{Offset: offset, Error: err}
		}
	}()

	return EventStore{
		Writer:           writer,
		WriterChan:       writerChan,
		WriterResultChan: writerResultChan,
	}
}

func (es *EventStore) WriteEvent(event domain.Event) (int64, error) {
	es.WriterChan <- event
	result := <-es.WriterResultChan
	return result.Offset, result.Error
}

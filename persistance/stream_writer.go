package persistance

import (
	"errors"
	"fmt"
	"os"

	"livenstore.evrard.online/persistance/serialization/writing"

	"github.com/oklog/ulid/v2"
)

var AlreadyInStreamError = errors.New("event already in stream")

type StreamWriter struct {
	BasePath string
}

func NewStreamWriter(basePath string) StreamWriter {
	return StreamWriter{
		BasePath: basePath,
	}
}

func (sw *StreamWriter) LinkEvent(streamName string, eventID ulid.ULID) error {
	alreadyInStream, err := sw.alreadyInStream(streamName, eventID)
	if err != nil {
		return err
	}
	if alreadyInStream {
		return AlreadyInStreamError
	}

	eventExists, err := sw.eventExists(eventID)
	if err != nil {
		return err
	}
	if !eventExists {
		return EventNotFound
	}

	payload := writing.WriteStringBytes(eventID.String(), []byte{})
	f, err := os.OpenFile(fmt.Sprintf("%s/%s.stream", sw.BasePath, streamName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(payload)
	if err != nil {
		return err
	}

	return nil
}

func (sw *StreamWriter) eventExists(eventID ulid.ULID) (bool, error) {
	reader := EventReader{BasePath: sw.BasePath}
	_, err := reader.ReadEvent(eventID)
	if err != nil {
		if err == EventNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (sw *StreamWriter) alreadyInStream(streamName string, eventID ulid.ULID) (bool, error) {
	reader := StreamReader{BasePath: sw.BasePath}
	stream, err := reader.ReadStream(streamName)
	if err != nil {
		return false, err
	}
	for _, id := range stream.EventIDs {
		if id == eventID {
			return true, nil
		}
	}
	return false, nil
}

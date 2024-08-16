package persistance

import (
	"fmt"
	"os"

	"github.com/oklog/ulid/v2"
	"livenstore.evrard.online/domain"
	"livenstore.evrard.online/persistance/serialization/reading"
)

type StreamReader struct {
	BasePath string
}

func NewStreamReader(basePath string) StreamReader {
	return StreamReader{
		BasePath: basePath,
	}
}

func (sr *StreamReader) ReadStream(streamName string) (*domain.Stream, error) {
	fileName := fmt.Sprintf("%s/%s.stream", sr.BasePath, streamName)

	stream := domain.Stream{
		Name: streamName,
	}

	f, err := os.OpenFile(fileName, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	b := make([]byte, stat.Size())
	_, err = f.Read(b)
	if err != nil {
		return nil, err
	}

	for len(b) > 0 {
		var eventID string
		eventID, b = reading.ReadStringBytes(b)
		ulid, err := ulid.Parse(eventID)
		if err != nil {
			return nil, err
		}

		stream.EventIDs = append(stream.EventIDs, ulid)
	}

	return &stream, nil
}

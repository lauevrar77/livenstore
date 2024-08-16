package persistance

import (
	"fmt"
	"os"

	"livenstore.evrard.online/persistance/serialization/writing"

	"github.com/oklog/ulid/v2"
)

type StreamWriter struct {
	BasePath string
}

func NewStreamWriter(basePath string) StreamWriter {
	return StreamWriter{
		BasePath: basePath,
	}
}

func (sw *StreamWriter) LinkEvent(streamName string, eventID ulid.ULID) error {
	//TODO : check that the event exists
	//TODO : check that the event is not already in the stream
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

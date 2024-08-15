package persistance

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/oklog/ulid/v2"
	"livenstore.evrard.online/domain"
	"livenstore.evrard.online/persistance/serialization/reading"
)

var EventNotFound = errors.New("event not found")

type EventReader struct {
	BasePath string
}

func (er *EventReader) ReadEvent(eventID ulid.ULID) (*domain.Event, error) {
	file, err := er.eventFile(eventID.String())
	if err != nil {
		return nil, err
	}

	f, err := os.Open(fmt.Sprintf("%s/%s", er.BasePath, file))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if stat.Size() == 0 {
		return nil, EventNotFound
	}

	b := make([]byte, stat.Size())
	_, err = f.Read(b)
	if err != nil {
		return nil, err
	}

	var e domain.Event
	for len(b) > 0 {
		e, b = reading.ReadEventBytes(b)
		if e.ID == eventID {
			return &e, nil
		}
	}

	return nil, EventNotFound
}

func (er *EventReader) eventFile(eventID string) (string, error) {
	files, err := os.ReadDir(er.BasePath)
	if err != nil {
		return "", err
	}
	eventFiles := []string{}
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".event") {
			continue
		}

		eventFiles = append(eventFiles, file.Name())
	}
	if len(eventFiles) == 0 {
		return "", EventNotFound
	}
	slices.Sort(eventFiles)

	for i, file := range eventFiles {
		firstEventID := strings.Split(file, ".")[0]
		if strings.Compare(firstEventID, eventID) >= 0 {
			if i == 0 {
				return "", EventNotFound
			}
			return eventFiles[i-1], nil
		}
	}
	return eventFiles[len(eventFiles)-1], nil
}

package persistance

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"livenstore.evrard.online/domain"
	"livenstore.evrard.online/persistance/serialization/writing"
)

const PAGE_SIZE = 1 * 1024 * 1024

var NOFILE = errors.New("No file found")

type EventWriter struct {
	BasePath    string
	currentFile *os.File
	currentSize int64
}

func NewEventWriter(basePath string) EventWriter {
	return EventWriter{
		BasePath: basePath,
	}
}

func (ew *EventWriter) WriteEvent(event domain.Event) (int64, error) {
	file, err := ew.CurrentFile(event.ID.String())
	if err != nil {
		return 0, err
	}

	pos, err := file.Seek(0, 1)
	if err != nil {
		return 0, err
	}

	eb := writing.WriteEventBytes(event, []byte{})
	_, err = file.Write(eb)
	if err != nil {
		return pos, err
	}

	ew.currentSize += int64(len(eb))

	return pos, nil
}

func (ew *EventWriter) CurrentFile(eventID string) (*os.File, error) {
	if ew.currentFile != nil {
		if ew.currentSize < PAGE_SIZE {
			return ew.currentFile, nil
		} else {
			ew.currentFile.Close()
		}
	}

	lastFile, err := ew.lastFile()
	if err != nil {
		if err == NOFILE {
			ew.currentSize = 0
			ew.currentFile, err = os.OpenFile(fmt.Sprintf("%s/%s.event", ew.BasePath, eventID), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
			if err != nil {
				return nil, err
			}
			return ew.currentFile, nil
		}
		return nil, err
	}

	ew.currentFile, err = os.OpenFile(lastFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return nil, err
	}
	stat, err := ew.currentFile.Stat()
	if err != nil {
		return nil, err
	}
	ew.currentSize = stat.Size()
	if ew.currentSize > PAGE_SIZE {
		ew.currentFile.Close()
		ew.currentSize = 0
		ew.currentFile, err = os.OpenFile(fmt.Sprintf("%s/%s.event", ew.BasePath, eventID), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			return nil, err
		}
	}
	return ew.currentFile, nil
}

func (ew *EventWriter) lastFile() (string, error) {
	files, err := os.ReadDir(ew.BasePath)
	if err != nil {
		return "", err
	}
	eventFiles := []string{}
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".event") {
			continue
		}

		eventFiles = append(eventFiles, fmt.Sprintf("%s/%s", ew.BasePath, file.Name()))
	}
	if len(eventFiles) == 0 {
		return "", NOFILE
	}
	slices.Sort(eventFiles)
	return eventFiles[len(eventFiles)-1], nil
}

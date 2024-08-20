package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/oklog/ulid/v2"
	"livenstore.evrard.online/domain"
	"livenstore.evrard.online/persistance"
	"livenstore.evrard.online/services"
	"livenstore.evrard.online/utils/bench"
)

var TEST_SIZE = 100_000
var NB_EVENTS = 1_000_000
var NB_STREAMS = 300

func main() {
	es := services.NewEventStore(
		"data",
		persistance.NewEventWriter,
		persistance.NewEventReader,
		persistance.NewStreamWriter,
		persistance.NewStreamReader,
	)

	fmt.Println("Creating events")
	eventIDs := CreateEvents(&es)

	fmt.Println("Benchmarking event by ID")
	times := make([]int64, TEST_SIZE)
	for i := 0; i < TEST_SIZE; i++ {
		before := time.Now().UnixMicro()
		es.EventByID(bench.SampleElement(eventIDs))
		after := time.Now().UnixMicro()
		times[i] = after - before
	}

	fmt.Println("Computing stats")
	timesJson, err := json.Marshal(bench.ComputeStats(times))
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("benchmarks/bench_event_by_id_1_client.json", timesJson, 0644)
}

func CreateEvents(es *services.EventStore) []ulid.ULID {
	eventIDs := make([]ulid.ULID, NB_EVENTS)
	streamNames := make([]string, NB_STREAMS)
	for i := 0; i < NB_STREAMS; i++ {
		streamNames[i] = bench.RandomString(10)
	}

	for i := 0; i < NB_EVENTS; i++ {
		eventID := ulid.Make()
		eventIDs[i] = eventID
		e := domain.Event{
			ID:        eventID,
			Type:      bench.SampleElement(streamNames),
			Timestamp: uint64(time.Now().Unix()),
			Data:      bench.RandomBytes(100),
		}
		es.PublishEvent(e)
	}

	return eventIDs
}

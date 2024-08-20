package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
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
var NB_CLIENTS = []int{2, 4, 8, 16}

func client(es *services.EventStore, eventIDs []ulid.ULID, outputChan chan int64, wg *sync.WaitGroup) {
	for i := 0; i < TEST_SIZE; i++ {
		before := time.Now().UnixMicro()
		es.EventByID(bench.SampleElement(eventIDs))
		after := time.Now().UnixMicro()
		outputChan <- (after - before)
	}
	close(outputChan)
	wg.Done()
}

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
	times := make(map[int][]int64)
	stats := make(map[int]bench.Stats)

	for _, nbClients := range NB_CLIENTS {
		fmt.Println("Running with", nbClients, "clients")

		var wg sync.WaitGroup
		outputChans := make([]chan int64, nbClients)
		for i := 0; i < nbClients; i++ {
			outputChan := make(chan int64, TEST_SIZE)
			outputChans[i] = outputChan
			go client(&es, eventIDs, outputChan, &wg)
			wg.Add(1)
		}

		wg.Wait()

		for _, outputChan := range outputChans {
			for t := range outputChan {
				times[nbClients] = append(times[nbClients], t)
			}
		}
		stats[nbClients] = bench.ComputeStats(times[nbClients])
		fmt.Println(stats[nbClients])
	}

	fmt.Println("Computing stats")
	timesJson, err := json.Marshal(stats)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("benchmarks/bench_event_by_id_multiple_client.json", timesJson, 0644)
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

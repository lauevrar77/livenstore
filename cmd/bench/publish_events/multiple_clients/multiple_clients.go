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

var PATH = "data"
var TEST_SIZE = 100_000
var NB_STREAMS = 30
var NB_CLIENTS = []int{2, 4, 8, 16, 32, 64, 128, 256}

func client(es *services.EventStore, outputChan chan int64, wg *sync.WaitGroup) {
	streamNames := make([]string, NB_STREAMS)
	for i := 0; i < NB_STREAMS; i++ {
		streamNames[i] = bench.RandomString(10)
	}
	for i := 0; i < TEST_SIZE; i++ {
		e := domain.Event{
			ID:        ulid.Make(),
			Type:      bench.SampleElement(streamNames),
			Timestamp: uint64(time.Now().Unix()),
			Data:      bench.RandomBytes(100),
		}
		before := time.Now().UnixMicro()
		es.PublishEvent(e)
		after := time.Now().UnixMicro()

		outputChan <- (after - before)
	}
	close(outputChan)
	wg.Done()
}

func wipeData(path string) {
	os.RemoveAll(path)
	os.Mkdir(path, 0755)
}

func main() {
	times := make(map[int][]int64)
	stats := make(map[int]bench.Stats)
	for _, nbClients := range NB_CLIENTS {
		fmt.Println("Running with", nbClients, "clients")
		wipeData(PATH)
		es := services.NewEventStore(
			PATH,
			persistance.NewEventWriter,
			persistance.NewEventReader,
			persistance.NewStreamWriter,
			persistance.NewStreamReader,
		)
		times[nbClients] = make([]int64, 0)

		wg := sync.WaitGroup{}
		clientChans := make([]chan int64, nbClients)
		for i := 0; i < nbClients; i++ {
			clientChans[i] = make(chan int64, TEST_SIZE)
			go client(&es, clientChans[i], &wg)
			wg.Add(1)
		}

		wg.Wait()

		fmt.Println("All clients are done. Aggregating results.")
		for _, clientChan := range clientChans {
			for time := range <-clientChan {
				times[nbClients] = append(times[nbClients], time)
			}
		}
		stats[nbClients] = bench.ComputeStats(times[nbClients])
		fmt.Println("Stats for", nbClients, "clients:", stats[nbClients])

		timesJson, err := json.Marshal(stats)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile("benchmarks/bench_publish_multiple_clients.json", timesJson, 0644)
	}
}

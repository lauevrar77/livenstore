package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/oklog/ulid/v2"
	"livenstore.evrard.online/domain"
	"livenstore.evrard.online/persistance"
	"livenstore.evrard.online/services"
	"livenstore.evrard.online/utils/bench"
)

var TEST_SIZE = 1_000_000

func main() {
	es := services.NewEventStore(
		"data",
		persistance.NewEventWriter,
		persistance.NewEventReader,
		persistance.NewStreamWriter,
		persistance.NewStreamReader,
	)

	times := make([]int64, TEST_SIZE)
	for i := 0; i < TEST_SIZE; i++ {
		e := domain.Event{
			ID:        ulid.Make(),
			Type:      bench.RandomString(10),
			Timestamp: uint64(time.Now().Unix()),
			Data:      bench.RandomBytes(100),
		}
		before := time.Now().UnixMicro()
		es.PublishEvent(e)
		after := time.Now().UnixMicro()

		times[i] = after - before
	}

	timesJson, err := json.Marshal(times)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("benchmarks/bench_publish_1_client.json", timesJson, 0644)
}

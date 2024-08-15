package main

import (
	"fmt"

	"github.com/oklog/ulid/v2"
	"livenstore.evrard.online/domain"
	"livenstore.evrard.online/persistance"
	"livenstore.evrard.online/services"
)

func main() {
	w := persistance.EventWriter{
		BasePath: "data",
	}
	r := persistance.EventReader{
		BasePath: "data",
	}
	es := services.NewEventStore(w)
	es.WriteEvent(domain.Event{
		ID:        ulid.Make(),
		Type:      "foo",
		Timestamp: 42,
		Data:      []byte{1, 2, 3},
	})

	ulid, err := ulid.Parse("01J5B9XDT2TY30R87D69ZWQ213")
	if err != nil {
		panic(err)
	}
	fmt.Println(r.ReadEvent(ulid))
}

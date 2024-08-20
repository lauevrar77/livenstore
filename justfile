run:
  go run cmd/server/server.go

bench_publish_1_client:
  rm data/*
  go run cmd/bench/publish_events/1_client/1_client.go

bench_publish_multiple_clients:
  go run cmd/bench/publish_events/multiples_clients/multiples_clients.go

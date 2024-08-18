run:
  go run cmd/server/server.go

bench_publish_1_client:
  rm data/*
  go run cmd/bench/publish_events/1_client.go

module github.com/rdev-org/api-gateway

go 1.27.1

require (
	github.com/go-chi/chi/v5 v5.3.2
	github.com/gorilla/websocket v1.5.3
	github.com/rdev-org/helpers v1.0.0
	github.com/rdev-org/proto-contracts v1.0.0
	github.com/segmentio/kafka-go v0.4.51
	go.yaml.in/yaml/v3 v3.0.5
	google.golang.org/grpc v1.84.0
)

require (
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260706201446-f0a921348800 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace github.com/rdev-org/proto-contracts => ../proto-contracts

replace github.com/rdev-org/helpers => ../helpers

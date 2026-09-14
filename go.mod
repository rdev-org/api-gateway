module github.com/rdev-org/api-gateway

go 1.27.1

require (
	github.com/go-chi/chi/v5 v5.3.2
	github.com/rdev-org/proto-contracts v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.83.2
)

require (
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260526163538-3dc84a4a5aaa // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

replace github.com/rdev-org/proto-contracts => ../proto-contracts

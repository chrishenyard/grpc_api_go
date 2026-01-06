# gRPC Jobs Service

A simple gRPC-based job management service built with Go and Protocol Buffers. This project demonstrates a basic CRUD API for managing job records using gRPC.

## Project Structure

```
grpc-jobs/
├── cmd/
│   ├── client/          # gRPC client implementation
│   └── server/          # gRPC server implementation
├── protos/              # Protocol Buffer definitions and generated code
│   ├── jobs.proto       # Service and message definitions
│   ├── jobs.pb.go       # Generated Go structs
│   └── jobs_grpc.pb.go  # Generated gRPC client/server code
└── .vscode/             # VS Code debug configuration
```

## Features

The service provides the following gRPC methods:

- **CreateJob**: Create a new job with name, description, and status
- **GetJob**: Retrieve a specific job by ID
- **GetJobs**: Stream all jobs with optional limit (server streaming)
- **DeleteJob**: Remove a job by ID

## API Definition

The service is defined in [`protos/jobs.proto`](protos/jobs.proto) with the following key message types:

- [`JobCreateRequest`](protos/jobs.pb.go): Contains job_name, job_description, and status
- [`JobRequest`](protos/jobs.pb.go): Contains job_id for lookup operations
- [`JobResponse`](protos/jobs.pb.go): Full job details including auto-generated job_id
- [`JobListOptions`](protos/jobs.pb.go): Options for listing jobs (limit parameter)

## Running the Service

### Prerequisites

- Go 1.19 or later
- Protocol Buffers compiler (protoc)
- gRPC Go plugins

### Start the Server

```bash
go run cmd/server/main.go
```

The server will start listening on port 50051.

### Run the Client

```bash
go run cmd/client/main.go
```

The client demonstrates creating a job, retrieving it, listing jobs, and deleting it.

## Implementation Details

### Server Implementation

The server ([`cmd/server/main.go`](cmd/server/main.go)) implements the [`JobsServer`](protos/jobs_grpc.pb.go) interface with:

- In-memory job storage using a map
- Thread-safe operations with `sync.RWMutex`
- UUID generation for job IDs using [`github.com/google/uuid`](https://github.com/google/uuid)
- Proper gRPC error handling with status codes

### Client Implementation

The client ([`cmd/client/main.go`](cmd/client/main.go)) demonstrates:

- Creating a gRPC connection with [`grpc.NewClient`](protos/jobs_grpc.pb.go)
- Using the generated [`JobsClient`](protos/jobs_grpc.pb.go) for method calls
- Handling streaming responses for the GetJobs method
- Context-based request timeouts

## Development

### Debugging

VS Code launch configurations are provided in [`.vscode/launch.json`](.vscode/launch.json) for debugging both client and server.

### Regenerating Protocol Buffers

If you modify [`protos/jobs.proto`](protos/jobs.proto), regenerate the Go code with:

```bash
protoc --go_out=. --go-grpc_out=. protos/jobs.proto
```

## Dependencies

- `google.golang.org/grpc` - gRPC Go implementation
- `google.golang.org/protobuf` - Protocol Buffers runtime
- `github.com/google/uuid` - UUID generation

## Error Handling

The service implements proper gRPC status codes:

- `InvalidArgument` for missing required fields
- `NotFound` for non-existent job IDs
- `Unimplemented` for unsupported operations (via [`UnimplementedJobsServer`](protos/jobs_grpc.pb.go))

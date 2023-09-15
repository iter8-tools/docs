# Get Iter8
FROM golang:1.21-bookworm as builder

# # Install protoc
# RUN apt update && apt install -y protobuf-compiler
# RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# RUN export export PATH="$PATH:$(go env GOPATH)/bin"

# Get source
WORKDIR /
COPY . ./
# COPY samples/abn-sample/frontend/abn.proto .

# # generate gRPC libraries from abn.proto
# RUN mkdir -p grpc
# RUN protoc --go_out=./grpc --go_opt=paths=source_relative --go-grpc_out=./grpc --go-grpc_opt=paths=source_relative  ./abn.proto

# Build frontend
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o fe main.go

### Multi-stage Docker build
### New image below

# Small Linux image with Iter8 binary
FROM scratch
WORKDIR /
COPY --from=builder /fe /frontend

ENTRYPOINT ["/frontend"]

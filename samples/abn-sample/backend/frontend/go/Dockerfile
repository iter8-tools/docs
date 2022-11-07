# Get Iter8
FROM golang:1.18-buster as builder

WORKDIR /
COPY frontend/go .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o fe main.go

### Multi-stage Docker build
### New image below

# Small linux image with iter8 binary
FROM debian:buster-slim
WORKDIR /
COPY --from=builder /fe /frontend

ENTRYPOINT ["/frontend"]

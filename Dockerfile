FROM golang:1.21 as build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd
COPY internal/ ./internal
COPY go.mod/ ./go.mod

RUN CGO_ENABLED=1 GOOS=linux go build -o /signals-collector cmd/main.go

ENTRYPOINT ["/signals-collector"]


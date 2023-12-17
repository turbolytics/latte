FROM golang:1.21 as build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd
COPY internal/ ./internal
COPY go.mod/ ./go.mod

RUN CGO_ENABLED=0 GOOS=linux go build -o /signals-collector cmd/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /signals-collector /signals-collector

ENTRYPOINT ["/signals-collector"]


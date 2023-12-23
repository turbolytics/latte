# signals-collector

Distill metrics from a variety of operational sources.

Signals Collector queries data at the source and only creates analytics data necessary for insights, leaving the operational data in the operational data store.

Signals collector produces analytic telemetry, allowing business analytics and business stakeholders to easily glean product insights from a variety of operational datasources including:

- Postgres
- Mongodb

<img width="600" alt="Screenshot 2023-12-23 at 7 45 56 AM" src="https://github.com/turbolytics/signals-collector/assets/151242797/15641829-7be3-4e73-b4b7-d7a9f1cc3f9b">

## Project Goals

Signals Collector aims to:

- Enable high fidelity data collection, creating the minimum dataset necessary to produce insights.
- Fidelity is focused informative signals and not raw data
- Aggregate data collection over raw data

## Approach
## Getting Started / Quickstart
## Examples

## Dev

- Start Postgres & Vector
```
docker-compose -f dev/compose.yaml up -d
```

### Invoke Locally

- Invoke collector
```
go run cmd/main.go config invoke --config=/Users/danielmican/code/github.com/turbolytics/signals-collector/dev/examples/postgres.http.stdout.yaml
```

- Verify Vector output
```
{"name":"core.users.signups.total","path":"/metrics","source_type":"http_server","tags":{"customer":"google"},"timestamp":"2023-12-12T13:33:27.564680676Z","type":"COUNT","value":3}
{"name":"core.users.signups.total","path":"/metrics","source_type":"http_server","tags":{"customer":"amazon"},"timestamp":"2023-12-12T13:33:27.564680676Z","type":"COUNT","value":2}
```

## Concepts
- Collector Configuration
- Source
- Sink
- Schedule
- Metric

## Docker

```
docker build -t signals/collector .
docker run -v $(PWD)/dev:/tmp signals-collector config validate --config=/tmp/examples/postgres.http.stdout.yaml
```

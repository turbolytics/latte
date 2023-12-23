# signals-collector
Collect metrics from a variety of operational sources

<img width="600" alt="Screenshot 2023-12-23 at 7 45 56 AM" src="https://github.com/turbolytics/signals-collector/assets/151242797/15641829-7be3-4e73-b4b7-d7a9f1cc3f9b">


## Project Goals
- High fidelity data collection, using the minimum data necessary
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

## Docker

```
docker build -t signals/collector .
docker run -v $(PWD)/dev:/tmp signals-collector config validate --config=/tmp/examples/postgres.http.stdout.yaml
```

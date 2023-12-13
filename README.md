# signals-collector
Collect metrics from a variety of operational sources

## Dev

- Start Postgres
- 
```
docker-compose -f dev/compose.yaml up -d
```

- Start Vector 
```
docker run -v $PWD/vector.yaml:/etc/vector/vector.yaml:ro -p 8686:8686 -p 9999:9999 timberio/vector:0.34.1-debian
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

## API

### Data Types

### REST

- `GET /collectors`
 
Lists all collectors.

- `POST /collectors`

Create a new collector.


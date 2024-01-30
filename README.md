# Latte 

<p align="center">
    <img width="400" alt="Screenshot 2024-01-21 at 9 59 19 AM" src="https://github.com/turbolytics/latte/assets/151242797/d675ed25-1459-4df3-94f7-e452d9962251">
</p>

Latte is a modern data engineering toolkit. Latte provides out of the box config-driven help to reduce the time and cost for many common data engineering tasks, including:

- Load
- Aggregate
- Transform
- Transfer 
- Extract 

Latte extracts and loads data from a variety of sources:

<img width="625" alt="Screenshot 2024-01-29 at 7 15 22 PM" src="https://github.com/turbolytics/latte/assets/151242797/1f00acda-933d-4786-babe-b52b7378b5b0">


## Project Goals

Latte aims to:

- Reduce cost of analytic data calculation, storage and query.
- Increase data warehouse data fidelity.
- Reduce data lake sprawl.
- Reduce data transfer into and costs of querying datalakes and datawarehouses.
- Enable product engineers to write self-service analytics built on their operational data stores.

Latte is configured using simple yaml configuration files:

<img width="911" alt="Screenshot 2023-12-24 at 9 05 22 AM" src="https://github.com/turbolytics/signals-collector/assets/151242797/80a59dcc-4f5e-4dd3-95bb-b9402cb3e6e7">


## Getting Started 

### Local Go Development

- Start docker dependencies
```
docker-compose -f dev/compose.yaml up -d
```

- Validate Configuration
```
go run cmd/main.go config validate --config=$(PWD)/dev/examples/postgres.kafka.stdout.yaml
VALID=true
```
- Invoke collector
```
go run cmd/main.go config invoke --config=$(PWD)/dev/examples/postgres.kafka.stdout.yaml
```

- Verify Kafka Output 
```
docker exec -it kafka1 kafka-console-consumer --bootstrap-server=localhost:9092 --topic=latte --from-beginning | jq .
```

```javascript
{
  "uuid": "4fdbc492-7e76-4ae7-9e32-82da7463f374",
  "name": "core.users.total",
  "value": 3,
  "type": "COUNT",
  "tags": {
    "customer": "google"
  },
  "timestamp": "2023-12-24T14:12:01.237538Z",
  "grain_datetime": "2023-12-24T00:00:00Z"
}
{
  "uuid": "cdc14916-15a4-4579-aa06-2dc65f442aba",
  "name": "core.users.total",
  "value": 2,
  "type": "COUNT",
  "tags": {
    "customer": "amazon"
  },
  "timestamp": "2023-12-24T14:12:01.237543Z",
  "grain_datetime": "2023-12-24T00:00:00Z"
}
```

### Local Docker 

- Start docker dependencies:
```
docker-compose -f dev/compose.with-collector.yaml up -d
```

- Validate config 
```
docker-compose -f dev/compose.with-collector.yaml run latte config validate --config=/dev/config/postgres.kafka.stdout.yaml

Creating dev_latte_run ... done
/dev/config/postgres.kafka.stdout.yaml
VALID=true
```

- Invoke Collector
```
docker-compose -f dev/compose.with-collector.yaml run latte config invoke --config=/dev/config/postgres.kafka.stdout.yaml

Creating dev_latte_run ... done
{"level":"info","ts":1703725219.4218154,"caller":"collector/collector.go:165","msg":"collector.Invoke","id":"d20161ac-be75-4868-8794-04d7bfa7d9d3","name":"postgres.users.total.24h"}
{"uuid":"a540fb6c-1638-4109-a385-3b0afda6fa12","name":"core.users.total","value":3,"type":"COUNT","tags":{"customer":"google"},"timestamp":"2023-12-28T01:00:19.422545549Z","grain_datetime":"2023-12-28T00:00:00Z"}
{"uuid":"fdac2b3f-a053-4997-b1a3-1ce1c6ca89a4","name":"core.users.total","value":2,"type":"COUNT","tags":{"customer":"amazon"},"timestamp":"2023-12-28T01:00:19.422548216Z","grain_datetime":"2023-12-28T00:00:00Z"}
```

- Verify Kafka Output
```
docker exec -it kafka1 kafka-console-consumer --bootstrap-server=localhost:9092 --topic=latte --from-beginning

{"uuid":"a540fb6c-1638-4109-a385-3b0afda6fa12","name":"core.users.total","value":3,"type":"COUNT","tags":{"customer":"google"},"timestamp":"2023-12-28T01:00:19.422545549Z","grain_datetime":"2023-12-28T00:00:00Z"}
{"uuid":"fdac2b3f-a053-4997-b1a3-1ce1c6ca89a4","name":"core.users.total","value":2,"type":"COUNT","tags":{"customer":"amazon"},"timestamp":"2023-12-28T01:00:19.422548216Z","grain_datetime":"2023-12-28T00:00:00Z"}
```

- Run Collector as daemon

```
docker-compose -f dev/compose.with-collector.yaml run latte run -c=/dev/config

Creating dev_latte_run ... done
{"level":"info","ts":1703726983.41746,"caller":"cmd/run.go:52","msg":"loading configs","path":"/dev/config"}
{"level":"info","ts":1703726983.4632287,"caller":"cmd/run.go:71","msg":"initialized collectors","num_collectors":5}
{"level":"info","ts":1703726983.4632988,"caller":"service/service.go:27","msg":"run"}
{"level":"info","ts":1703726983.4633727,"caller":"collector/collector.go:165","msg":"collector.Invoke","id":"5ba44984-a8a3-42ac-a70d-85e11f808a6c","name":"postgres.users.total.24h"}
```

- Tail the audit log, Check Kafka, Verify Vector 
 
```
tail -f dev/audit/latte.audit.log
```

## Examples

### Example Configurations

Checkout the [examples directory](./dev/examples) for configuration examples.

## Additional Documentation

Additional documentation is available in the [docs/ directory](./docs)

----

<a href="https://www.freepik.com/free-vector/cute-happy-coffee-cup-cartoon-vector-icon-illustration-drink-character-icon-concept-flat-cartoon-style_12090270.htm#query=coffee%20cartoon&position=14&from_view=search&track=ais&uuid=488418ce-32c1-4adf-9001-a7409e29a639">Image by catalyststuff</a> on Freepik
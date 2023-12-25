# signals-collector

Signals collector is a modern Extract Transform and Aggregate (ETA) tool. 

Signals Collector queries data at the source and only creates analytics data necessary for insights, leaving the operational data in the operational data store.

Signals collector produces analytic telemetry, allowing business analytics and business stakeholders to easily glean product insights from a variety of operational datasources including:

- Postgres
- Mongodb

<img width="600" alt="Screenshot 2023-12-23 at 7 45 56 AM" src="https://github.com/turbolytics/signals-collector/assets/151242797/15641829-7be3-4e73-b4b7-d7a9f1cc3f9b">

Signals Collector is a new class of tooling for Extract Transform and Aggregation (ETA). Signals collector aggregates data at the source and only emits the aggregations for downstream processing in the datalake or datawarehouse:

<img width="600" alt="Screenshot 2023-12-24 at 9 02 52 AM" src="https://github.com/turbolytics/signals-collector/assets/151242797/f32cf84c-e05f-4a59-8040-8f3cc74b04a6">

Consider a users table with hundreds of thousands or millions of users:

```
CREATE TABLE users (
    account varchar,
    signup_time timestamp
);
```

Signals collector aggregates data at the source and only emits the aggregated data, at the grain consumed by end users. Imagine the users has 2 customers, amazon and google, with 1MM users:


Signals collector can aggregate user counts at the source:

<img width="911" alt="Screenshot 2023-12-24 at 9 05 22 AM" src="https://github.com/turbolytics/signals-collector/assets/151242797/80a59dcc-4f5e-4dd3-95bb-b9402cb3e6e7">


This will produce 2 data points per day, reducing the need to egress millions of users to a datalake or data warehouse:

```
{
  "uuid": "4fdbc492-7e76-4ae7-9e32-82da7463f374",
  "name": "core.users.total",
  "value": 750000,
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
  "value": 250000,
  "type": "COUNT",
  "tags": {
    "customer": "amazon"
  },
  "timestamp": "2023-12-24T14:12:01.237543Z",
  "grain_datetime": "2023-12-24T00:00:00Z"
}
```

Signals Collector takes a different approach to data analytics when compared to tools such as [fivetran](https://www.fivetran.com/) and [airbyte](https://airbyte.com/). These tools focus on copying operational data sources for downstream analysis:

<img width="600" alt="Screenshot 2023-12-23 at 7 45 56 AM" src="https://github.com/turbolytics/signals-collector/assets/151242797/d17f07ef-5744-4210-a652-f836ceb399df">

These tools create a 1:1 copy of the operational product data in the datalake or datawarehouse. Once in a datalake or datawarehouse the operational data requires layers and layers of processing to distill insights:

<img width="600" alt="Screenshot 2023-12-23 at 7 45 56 AM" src="https://github.com/turbolytics/signals-collector/assets/151242797/1b3df140-df6d-4d71-a47a-f36e38986a20">

The operational data is often at a very fine grain, individual user product transactions. Data tools are used to slowly refine and aggregate the data into a coarse grain consumable by humans. This refinement is extremely time-consuming and flaky:

<img width="600" alt="Screenshot 2023-12-23 at 7 45 56 AM" src="https://github.com/turbolytics/signals-collector/assets/151242797/08185885-10fa-4c3f-8df7-fa29678fb5f2">

Some use cases such as data exploration and machine learning do benefit from large sets of operational data, but business analytics often doesn't. The layered warehouse approach is testament to this. The final layer of data is aggregated (often to the day) and very small compared to the fine-grained large source data. The modern data stack encourages egressing of huge amounts of operational data, only to aggregate that data over a series of stages into a coarse grained aggregate data set!!!

https://on-systems.tech/blog/135-draining-the-data-swamp/

Signals Collector enables aggregating analytic data at the source and only emit the data that is necessary for analytics at the grain necessary. 

## Project Goals

Signals Collector aims to:

- Reduce cost of analytic data calculation, storage and query.
- Increase data warehouse data fidelity.
- Reduce data lake sprawl.
- Reduce data transfer into and costs of querying datalakes and datawarehouses.
- Enable product engineers to write self service analytics built on their operational data stores.


## Getting Started / Quickstart

- Start docker dependencies
```
docker-compose -f dev/compose.yaml up -d
```

### Invoke Locally

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
docker exec -it kafka1 kafka-console-consumer --bootstrap-server=localhost:9092 --topic=signals --from-beginning | jq .
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
## Examples

### Example Configurations

Checkout the [examples directory](./dev/examples) for configuration examples.

## Concepts
### Collector Configuration

### Source
### Sink
### Schedule
### Metric

## Docker

```
docker build -t signals/collector .
docker run -v $(PWD)/dev:/tmp signals-collector config validate --config=/tmp/examples/postgres.http.stdout.yaml
```


## Additional Documentation

Additional documention is available in the [docs/ directory](./docs)
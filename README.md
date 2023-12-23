# signals-collector

Distill metrics from a variety of operational sources.

Signals Collector queries data at the source and only creates analytics data necessary for insights, leaving the operational data in the operational data store.

Signals collector produces analytic telemetry, allowing business analytics and business stakeholders to easily glean product insights from a variety of operational datasources including:

- Postgres
- Mongodb

<img width="600" alt="Screenshot 2023-12-23 at 7 45 56 AM" src="https://github.com/turbolytics/signals-collector/assets/151242797/15641829-7be3-4e73-b4b7-d7a9f1cc3f9b">

Signals Collector takes a different approach to data analytics when compared to tools such as [fivetran](https://www.fivetran.com/) and [airbyte](https://airbyte.com/). These tools focus on copying operational data sources for downstream analysis:

<img width="600" alt="Screenshot 2023-12-23 at 7 45 56 AM" src="https://github.com/turbolytics/signals-collector/assets/151242797/d17f07ef-5744-4210-a652-f836ceb399df">

These tools create a 1:1 copy of the operational product data in the datalake or datawarehouse. Once in a datalake or datawarehouse the operational data requires layers and layers of processing to distill insights:

<img width="600" alt="Screenshot 2023-12-23 at 7 45 56 AM" src="https://github.com/turbolytics/signals-collector/assets/151242797/1b3df140-df6d-4d71-a47a-f36e38986a20">

The operational data is often at a very fine grain, individual user product transactions. Data tools are used to slowly refine and aggregate the data into a coarse grain consumable by humans. This refinement is extremely time-consuming and flaky:

<img width="600" alt="Screenshot 2023-12-23 at 7 45 56 AM" src="https://github.com/turbolytics/signals-collector/assets/151242797/08185885-10fa-4c3f-8df7-fa29678fb5f2">

Some use cases such as data exploration and machine learning do benefit from large sets of operational data, but business analytics often doesn't. The layered warehouse approach is testament to this. The final layer of data is aggregated (often to the day) and very small compared to the fine-grained large source data. The modern data stack encourages egressing of huge amounts of operational data, only to aggregate that data over a series of stages into a coarse grained aggregate data set!!!

https://on-systems.tech/blog/135-draining-the-data-swamp/

Signals Collector enables aggregating analytic data at the source and only emit the data that is necessary for analytics at the grain necessary. Aggregating at the source has a number of benefits including:

- Reduce cost of analytic data calculation, storage and query.
- Increase datawarehouse data fidelity.
- Reduce data lake sprawl.
- Reduce data transfer into and costs of querying datalakes and datawarehouses.
- Enable product engineers to write self service analytics built on their operational data stores.

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

# Prometheus SLO Storage in a System of Record

Prometheus is a timeseries data store for operational telemetry. Prometheus is used to determine health of operational systems, debug issues and alert on issues when they do occur. Many companies have the need to report on operational telemetry in the form of Service Level Objectives (SLOs). SLOs determine service availability. 

This tutorial shows how signals collector is used to capture prometheus SLO data for long term storage in a datalake or datawarehouse.

<img width="713" alt="Screenshot 2024-01-07 at 1 01 52 PM" src="https://github.com/turbolytics/signals-collector/assets/151242797/c8d9f6fe-8a31-42c9-85a1-00001f2fb3c4">

The signals collector will:

1. Scrape prometheus API daily for yesterday's data 
2. Transform the prometheus API response using duckdb sql 
3. Sink the aggregate data to the configured datastores (kakfa, file audit, vector/HTTP, etc)

## Configuration File

<img width="927" alt="Screenshot 2024-01-04 at 4 58 44 PM" src="https://github.com/turbolytics/signals-collector/assets/151242797/64423f0f-360c-4fe3-8a86-16251a59de6b">

## Tutorial

### Determine the Prometheus Query

The first step is to determine the prometheus query that will produce the correct daily aggregate data. Prometheus ships with a web frontend that supports querying underlying prometheus data.

<img width="1512" alt="Screenshot 2024-01-15 at 2 18 45 PM" src="https://github.com/turbolytics/signals-collector/assets/151242797/51f3e14a-8aa8-4816-9f2c-6cde43fd8756">

[Test the query directly against the API to see the response returned](https://prometheus.io/docs/prometheus/latest/querying/api/#instant-queries):

```
curl 'http://localhost:9090/api/v1/query?query=increase%28collector_invoke_count_total%5B24h%5D%29' | jq .
```

```
{
  "status": "success",
  "data": {
    "resultType": "vector",
    "result": [
      {
        "metric": {
          "collector_name": "mongo.users.10s",
          "instance": "host.docker.internal:12223",
          "job": "signals-collector-host",
          "otel_scope_name": "signals-collector",
          "result_status_code": "OK"
        },
        "value": [
          1704403299.513,
          "65.28921807134175"
        ]
      },
    ...
```

### Configure Signals Collector to Query Prometheus API

Signals collector needs to know where to reach prometheus and the query to execute against prometheus:

```yaml
source:
  type: prometheus
  config:
    uri: 'http://{{ getEnvOrDefault "SC_PROM_HOST" "localhost" }}:9090/api/v1/query'
    query: 'increase(collector_invoke_count_total[24h])'
    time_expression: $start_of_day
    ...
```

Notice the `time_expression`. This is a special variable that will pass in `&time=<UNIX_TIMESTAMP>` parameter into the `query`. This will ensure idempotent data collection runs. Every time signals collector runs it will start collection from the beginning of the day. This will ensure that the date is consistent, and that the date being used is already passed, meaning that the full 24 hours of availability data is present.

### Configure Signals Collector Transform 

The next step is to specify the transform that will be applied to the prometheus API output:

```yaml
source:
  type: prometheus
  config:
    ...
    sql: |
      SELECT
        metric->>'collector_name' as service,
        CASE
          WHEN metric->>'result_status_code' = 'OK'
          THEN false
          ELSE true
        END as error,
        round(value::DOUBLE, 0) as value
      FROM
        prom_metrics
```

This sql is executed using duckdb against the API response returned from the prometheus API. Using duckdb provides the end user the power to model and transform the responses using SQL. Output from this stage will look similiar to the following:

```json
...
{
  "uuid": "6fc67572-76d1-4c36-a9d1-e07a5a40f5cb",
  "name": "service.availability.24h",
  "value": 12,
  "type": "COUNT",
  "tags": {
    "env": "prod",
    "error": "false",
    "service": "postgres.users.total.1m"
  },
  "timestamp": "2024-01-04T21:30:50.345872Z",
  "grain_datetime": "2024-01-04T00:00:00Z"
}
...
```

### Validate Signals Collector Config

The next step is to validate the configuration:

```
go run cmd/main.go config validate --config=$(PWD)/dev/examples/prometheus.stdout.yaml

VALID=true
```

### Invoke Signals Collector Config

After the collector configuration is validated it can be invoked, using the following command:

```
go run cmd/main.go config invoke --config=$(PWD)/dev/examples/prometheus.stdout.yaml
```

### Start Signals collector for regular collection

```
go run cmd/main.go run -c /path/to/your/configuration/file.yaml
```

### Query Output for Reporting

The final step is to query data for reporting. This example collects availability data daily. Error and non-error counts are collected which allow for availability calculation over arbitrary intervals 

sum(success) 
/ 
sum(successes + errors)



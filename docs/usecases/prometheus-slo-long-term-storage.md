# Prometheus SLO Storage in a System of Record

Prometheus is a timeseries data store for operational telemetry. Prometheus is used to determine health of operational systems, debug issues and alert on issues when they do occur. Many companies have the need to report on operational telemetry in the form of Service Level Objectives (SLOs). SLOs determine service availability. 

This tutorial shows how latte is used to capture prometheus SLO data for long term storage in a datalake or datawarehouse.

<img width="676" alt="Screenshot 2024-01-29 at 7 11 53 PM" src="https://github.com/turbolytics/latte/assets/151242797/4023f04e-9bf3-4538-ae5e-fee750bc7350">

The latte will:

1. Scrape prometheus API daily for yesterday's data 
2. Transform the prometheus API response using duckdb sql 
3. Sink the aggregate data to the configured datastores (kakfa, file audit, vector/HTTP, etc)

## Configuration File

<img width="887" alt="Screenshot 2024-01-15 at 2 32 32 PM" src="https://github.com/turbolytics/signals-collector/assets/151242797/c3beb697-7b86-4624-8f60-688216c4ad01">

## Tutorial

### Determine the Prometheus Query

The first step is to determine the prometheus query that will produce the correct daily aggregate data. Prometheus ships with a web frontend that supports querying underlying prometheus data:

<img width="1512" alt="Screenshot 2024-01-15 at 2 18 45 PM" src="https://github.com/turbolytics/signals-collector/assets/151242797/51f3e14a-8aa8-4816-9f2c-6cde43fd8756">

Once a query is built, [test the query directly against the API to see the response returned](https://prometheus.io/docs/prometheus/latest/querying/api/#instant-queries):

```
curl 'http://localhost:9090/api/v1/query?query=sum%28round%28increase%28collector_invoke_count_total%5B24h%5D%29%29%29%20by%20%28result_status_code%29' | jq .
```

```json
{
  "status": "success",
  "data": {
    "resultType": "vector",
    "result": [
      {
        "metric": {
          "result_status_code": "OK"
        },
        "value": [
          1705347547.606,
          "23"
        ]
      },
      {
        "metric": {
          "result_status_code": "ERROR"
        },
        "value": [
          1705347547.606,
          "0"
        ]
      }
    ]
  }
}
```

### Configure latte to Query Prometheus API

latte needs to know where to reach prometheus and the query to execute against prometheus:

```yaml
  type: prometheus
  strategy: historic_tumbling_window 
  config:
    uri: 'http://{{ getEnvOrDefault "SC_PROM_HOST" "localhost" }}:9090/api/v1/query'
    query: 'sum(round(increase(collector_invoke_count_total[24h]))) by (result_status_code)'
    time:
      window: 24h
```

The `type` value specifies that this is a `prometheus` collection. The strategy specifies that this is a `window` type collection. More information about `window` collection strategy is available here. 

Notice the `time.window`. This is a special variable that tells prometheus how much data it is working with. The collector needs this information to correctly schedule data collections against prometheus. The `time.window` property is used to avoid parsing this information out of the prometheus `query` value.

This will ensure idempotent data collection runs. Every time latte runs it will start collection from the beginning of the day. This will ensure that the date is consistent, and that the date being used is already passed, meaning that the full 24 hours of availability data is present.

### Configure latte Transform 

The next step is to specify the transform that will be applied to the prometheus API output:

```yaml
source:
  type: prometheus
  config:
    ...
    sql: |
      SELECT
          'latte' as service,
          CASE
              WHEN metric->>'result_status_code' = 'OK'
              THEN false
              ELSE true
          END as error,
          round(value::DOUBLE, 0) as value
      FROM
          prom_metrics
```

This sql is executed using duckdb against the API response returned from the prometheus API. Using duckdb provides the end user the power to model and transform the responses using SQL. Output from this stage will look similar to the following:

```json
{
  "uuid": "fa201e5e-c67c-451c-a321-f51924a08d57",
  "name": "service.availability.24h",
  "value": 510,
  "type": "COUNT",
  "tags": {
    "env": "prod",
    "error": "false",
    "service": "latte"
  },
  "timestamp": "2024-01-15T19:11:31.637377Z",
  "window": "2024-01-15T00:00:00Z"
}
```

### Validate latte Config

The next step is to validate the configuration:

```
go run cmd/main.go config validate --config=$(PWD)/dev/examples/prometheus.stdout.yaml

VALID=true
```

Validation ensures that the configuration file is correct and can be successfully parsed.

### Invoke latte Config

After the collector configuration file is validated it can be invoked using the following command:

```
go run cmd/main.go config invoke --config=$(PWD)/dev/examples/prometheus.stdout.yaml
```

Invocation actually executes the configuration file. This is used to verify that data is sourced, transformed and sunk correctly.

### Start latte for Regular Collection

Signal collector ships with a daemon that executes collectors at their specified intervals. The following command starts latte as a service: 

```
go run cmd/main.go run -c $(PWD)/dev/examples/prometheus.stdout.yaml
```

The service executes the prometheus collector every 12 hours. The prometheus collector will collect availability data for 24 hours and emit a single data point every time it encounters a fully complete day. More information about this is available under the "window" strategy documentation.

latte also ships as a [docker image](https://github.com/turbolytics/latte/blob/main/Dockerfile) and on [dockerhub](https://hub.docker.com/repository/docker/turbolytics/latte/general).

### Query Output for Reporting

The final step is to query data for reporting. This example collects availability data daily. Error and non-error counts are collected which allow for availability calculation over arbitrary intervals. The prometheus collector used in this example outputs availibility data to a local audit file. The following shows an example of this data:

```
D select name, value, tags, "window" from read_json('latte.audit.log', auto_detect=true, format=newline_delimited);
┌──────────────────────────┬───────┬─────────────────────────────────────────────────────────────┬─────────────────────┐
│           name           │ value │                            tags                             │       window        │
│         varchar          │ int64 │     struct(env varchar, error varchar, service varchar)     │      timestamp      │
├──────────────────────────┼───────┼─────────────────────────────────────────────────────────────┼─────────────────────┤
│ service.availability.24h │    41 │ {'env': prod, 'error': false, 'service': latte}             │ 2024-01-15 00:00:00 │
│ service.availability.24h │    61 │ {'env': prod, 'error': false, 'service': latte}             │ 2024-01-14 00:00:00 │
│ service.availability.24h │     2 │ {'env': prod, 'error': true, 'service': latte}              │ 2024-01-14 00:00:00 │
│ service.availability.24h │    80 │ {'env': prod, 'error': false, 'service': latte}             │ 2024-01-13 00:00:00 │
│ service.availability.24h │    71 │ {'env': prod, 'error': false, 'service': latte}             │ 2024-01-12 00:00:00 │
└──────────────────────────┴───────┴─────────────────────────────────────────────────────────────┴─────────────────────┘
```

The collector counts successes and errors per day. Notice how some collections do not include errors.

Availability can be calculated after the data is collected. Availability is a ratio and calculated over an interval. Availability is defined as:

```
sum(success) 
/ 
sum(successes + errors)
```

```
SELECT
    "window",
    ROUND((sum(successes) / (sum(successes) + sum(errors))), 4) as availability
FROM (
    SELECT 
      "window",
      CASE
        WHEN tags ->> 'error' = true 
          THEN value
          ELSE 0
        END AS errors,
      CASE
        WHEN tags ->> 'error' = false 
          THEN value
            ELSE 0
        END AS successes
    FROM   
      read_json('latte.audit.log', auto_detect = true, format = newline_delimited)
    GROUP  BY 
      "window",
      tags,
      value
    ORDER BY 
      "window" DESC  
)
GROUP BY 
    "window"
ORDER BY 
    "window"
```

The final result is availability data by day!
```
┌─────────────────────┬──────────────┐
│       window        │ availability │
│      timestamp      │    double    │
├─────────────────────┼──────────────┤
│ 2024-01-12 00:00:00 │          1.0 │
│ 2024-01-13 00:00:00 │          1.0 │
│ 2024-01-14 00:00:00 │       0.9683 │
│ 2024-01-15 00:00:00 │          1.0 │
└─────────────────────┴──────────────┘
```



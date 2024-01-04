# Prometheus SLO Storage in a System of Record

Prometheus is a timeseries data store for operational telemetry. Prometheus is used to determine health of operational systems, debug issues and alert on issues when they do occur. Many companies have the need to report on operational telemetry in the form of Service Level Objectives (SLOs). SLOs determine service availability. 

This tutorial shows how signals collector is used to capture prometheus SLO data for long term storage in a datalake or datawarehouse.

<diagram>

The signals collector will:
- Scrape prometheus API daily for daily SLO data, 
- Transform the prometheus API response using duckdb sql 
- Sink the aggregate data to the configured datastores (kakfa, file audit, vector/HTTP, etc)

## Configuration File

## Tutorial

### Determine the Prometheus Query

The first step is to determine the prometheus query that will produce the correct daily aggregate data.



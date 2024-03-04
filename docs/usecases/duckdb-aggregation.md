

```
$ awslocal s3api create-bucket --bucket localstack-latte
$ awslocal s3 cp access.json s3://localstack-latte/activity-events/year=2024/month=3/day=3/hour=1/
$ go run cmd/main.go config invoke --config=$(PWD)/dev/examples/s3.kafka.yaml --start-time='2024-03-03T00:00:00Z'
```


test-unit:
	go test -short ./...

test-integration:
	go test -run TestIntegration ./...

docker-image:
	docker build -t turbolytics/latte .

test: test-unit test-integration

.PHONY: test-unit test-integration docker-image
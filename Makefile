

test-unit:
	go test -short ./...

test-integration:
	go test -run TestIntegration ./...


test: test-unit test-integration

.PHONY: test-unit test-integration
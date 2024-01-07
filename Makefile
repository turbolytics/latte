

test-unit:
	go test -short ./... -v

test-integration:
	go test -run TestIntegration ./... -v


test: test-unit test-integration

.PHONY: test-unit test-integration
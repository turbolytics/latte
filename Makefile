

test-unit:
	go test -short ./... -v

test-integration:
	go test -run TestIntegration ./... -v

.PHONY: test-unit test-integration
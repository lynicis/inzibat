.PHONY: generate-mock
generate-mock:
	mockgen -source=config/reader.go -destination=config/reader_mock.go -package=config

.PHONY: lint
lint:
	golangci-lint --verbose run ./...

.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

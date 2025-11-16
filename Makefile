.PHONY: generate-mock
generate-mock:
	mockgen -source=config/reader.go -destination=config/reader_mock.go -package=config
	mockgen -source=cmd/form_builder/form_collectors.go -destination=cmd/form_builder/form_runner_mock.go -package=form_builder
	mockgen -source=cmd/create.go -destination=cmd/create_mock.go -package=cmd

.PHONY: lint
lint:
	golangci-lint --verbose run ./...

.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

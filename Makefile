.PHONY: generate-mock lint coverage sonar-up sonar-down sonar-scan vuln sec

mockgen = "go tool go.uber.org/mock/mockgen"
golangci-lint = "go tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint"
govulncheck = "go tool golang.org/x/vuln/cmd/govulncheck"
gosec = "go tool github.com/securego/gosec/v2/cmd/gosec"

generate-mock:
	mockgen -source=config/reader.go -destination=config/reader_mock.go -package=config
	mockgen -source=cmd/formbuilder/form_collector.go -destination=cmd/formbuilder/form_runner_mock.go -package=formbuilder

lint:
	golangci-lint --verbose run ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

vuln:
	govulncheck ./...

sec:
	go tool gosec ./...

sonar-up:
	@echo "Starting local SonarQube container..."
	docker run -d --name sonarqube -e SONAR_ES_BOOTSTRAP_CHECKS_DISABLE=true -p 9000:9000 sonarqube:lts-community
	@echo "Waiting for SonarQube to be up... (this may take a minute)"
	@while ! curl -s http://localhost:9000/api/system/status | grep -q '"status":"UP"'; do \
		sleep 2; \
	done
	@echo "SonarQube is up! Changing default admin password to admin123..."
	@curl -s -u admin:admin -X POST "http://localhost:9000/api/users/change_password?login=admin&previousPassword=admin&password=admin123" >/dev/null 2>&1 || true

sonar-down:
	@echo "Stopping and removing local SonarQube container..."
	docker stop sonarqube || true
	docker rm sonarqube || true

sonar-scan: coverage
	@echo "Running SonarScanner via Docker..."
	docker run --rm \
		-v "$(PWD):/usr/src" \
		-e SONAR_HOST_URL="http://host.docker.internal:9000" \
		sonarsource/sonar-scanner-cli \
		-Dsonar.projectKey=inzibat \
		-Dsonar.sources=. \
		-Dsonar.exclusions="**/*_mock.go" \
		-Dsonar.go.coverage.reportPaths=coverage.out \
		-Dsonar.login=admin \
		-Dsonar.password=admin123 \
		-Dsonar.branch.name=




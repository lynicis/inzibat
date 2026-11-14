.PHONY: generate-mock
generate-mock:
	mockgen -source=config/reader.go -destination=config/reader_mock.go -package=config
	mockgen -source=cmd/formbuilder/form_collector.go -destination=cmd/formbuilder/form_runner_mock.go -package=formbuilder
	mockgen -source=cmd/create.go -destination=cmd/create_mock.go -package=cmd

.PHONY: lint
lint:
	golangci-lint --verbose run ./...

.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: sonar-up
sonar-up:
	@echo "Starting local SonarQube container..."
	docker run -d --name sonarqube -e SONAR_ES_BOOTSTRAP_CHECKS_DISABLE=true -p 9000:9000 sonarqube:lts-community
	@echo "Waiting for SonarQube to be up... (this may take a minute)"
	@while ! curl -s http://localhost:9000/api/system/status | grep -q '"status":"UP"'; do \
		sleep 2; \
	done
	@echo "SonarQube is up! Changing default admin password to admin123..."
	@curl -s -u admin:admin -X POST "http://localhost:9000/api/users/change_password?login=admin&previousPassword=admin&password=admin123" >/dev/null 2>&1 || true

.PHONY: sonar-down
sonar-down:
	@echo "Stopping and removing local SonarQube container..."
	docker stop sonarqube || true
	docker rm sonarqube || true

.PHONY: sonar-scan
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




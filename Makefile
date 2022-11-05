.PHONY: infra-cleanup
infra-cleanup:
	@docker-compose down --volumes
	@docker rmi rest-articles_rest
	@docker container prune -f

# 1. Database Related
.PHONY: database-unit
database-unit:
	@docker-compose --profile unittest up --build

# 2. Unit Testing
.PHONY: test
test: ## Run tests with check race and coverage
	@go test -failfast -count=1 ./... -json -cover -race | tparse -smallscreen

# 3. Run
.PHONY: run
run:
	@docker-compose --profile run up --build
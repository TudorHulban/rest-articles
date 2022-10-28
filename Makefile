.PHONY: infra
infra:
	@make infra-cleanup
	@docker-compose up --renew-anon-volumes

.PHONY: infra-cleanup
infra-cleanup:
	@docker-compose down --volumes

# 1. Database Related
.PHONY: database-unit
database-unit:
	@docker-compose  --profile unittest  up

# 2. unit Testing
.PHONY: test
test: ## Run tests with check race and coverage
	@go test -failfast -count=1 ./... -json -cover -race | tparse -smallscreen
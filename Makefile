.PHONY: infra
infra:
	@make infra-cleanup
	@docker-compose up --renew-anon-volumes

.PHONY: infra-cleanup
infra-cleanup:
	@docker-compose down --volumes

.PHONY: test
test: ## Run tests with check race and coverage
	@go test -failfast -count=1 ./... -json -cover -race | tparse -smallscreen

migrate-up:
	migrate -path db/migration -database ${db} -verbose up

migrate-down:
	migrate -path db/migration -database ${db} -verbose down
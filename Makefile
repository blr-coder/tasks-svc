GOLANG_MIGRATE_VERSION := v4.17.1
GOLANG_CI_LINT_VERSION := v1.59.1
GCI_VERSION := v0.10.1
GO_MOCK_VERSION := v1.6.0

MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
BUILD_PATH := $(dir $(MAKEFILE_PATH))
GOBIN ?= $(BUILD_PATH)tools/bin

MIGRATION_DIR = internal/infrastructure/storages/psql_store/migrations
MIGRATION_URL = "postgres://task_svc_db_user:task_svc_db_user_pass@localhost:5232/task_svc_db?sslmode=disable"

migrate_new: # pass name as parameter, example - make migrate_new name=add_task_table
	$(GOBIN)/golang-migrate/$(GOLANG_MIGRATE_VERSION)/migrate create -ext sql -dir $(MIGRATION_DIR)/ -seq $(name)
migrate_up:
	$(GOBIN)/golang-migrate/$(GOLANG_MIGRATE_VERSION)/migrate -path $(MIGRATION_DIR) -database $(MIGRATION_URL) -verbose up
migrate_down:
	$(GOBIN)/golang-migrate/$(GOLANG_MIGRATE_VERSION)/migrate -path $(MIGRATION_DIR) -database $(MIGRATION_URL) -verbose down 1
migrate_status:
	$(GOBIN)/golang-migrate/$(GOLANG_MIGRATE_VERSION)/migrate -path $(MIGRATION_DIR) -database $(MIGRATION_URL) version


.PHONY: test
test:
	 go test ./...

.PHONY: lint
lint:
	$(GOBIN)/golangci-lint/$(GOLANG_CI_LINT_VERSION)/golangci-lint run --config .golangci.yml

.PHONY: generate
generate:
	@MOCKGEN_PATH="$(GOBIN)/mockgen/$(GO_MOCK_VERSION)/mockgen" $(goflags) go generate ./...

.PHONY: tools
tools:
	@if [ ! -f $(GOBIN)/golangci-lint ]; then\
		echo "Installing golangci-lint $(GOLANG_CI_LINT_VERSION)";\
		GOBIN=$(GOBIN)/golangci-lint/$(GOLANG_CI_LINT_VERSION) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANG_CI_LINT_VERSION);\
		echo "Done";\
	fi

	@if [ ! -f $(GOBIN)/gci ]; then\
		echo "Installing gci $(GCI_VERSION)";\
		GOBIN=$(GOBIN)/gci/$(GCI_VERSION) go install github.com/daixiang0/gci@$(GCI_VERSION);\
		echo "Done";\
	fi

	@if [ ! -f $(GOBIN)/golang-migrate ]; then\
    		echo "Installing golang-migrate $(GOLANG_MIGRATE_VERSION)";\
    		GOBIN=$(GOBIN)/golang-migrate/$(GOLANG_MIGRATE_VERSION) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$(GOLANG_MIGRATE_VERSION);\
    		echo "Done";\
    fi

	@if [ ! -f $(GOBIN)/mockgen ]; then\
		echo "Installing mockgen $(GO_MOCK_VERSION)";\
		GOBIN=$(GOBIN)/mockgen/$(GO_MOCK_VERSION) go install github.com/golang/mock/mockgen@$(GO_MOCK_VERSION);\
		echo "Done";\
	fi

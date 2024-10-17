MIGRATION_DIR = internal/infrastructure/storages/psql_store/migrations
MIGRATION_URL = "postgres://task_svc_db_user:task_svc_db_user_pass@localhost:5232/task_svc_db?sslmode=disable"

MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
BUILD_PATH := $(dir $(MAKEFILE_PATH))
GOBIN ?= $(BUILD_PATH)tools/bin

GOLANG_CI_VERSION := v1.59.1
GCI_VERSION := v0.10.1

migrate_new: # pass name as parameter, example - make migrate_new name=add_task_table
	migrate create -ext sql -dir $(MIGRATION_DIR)/ -seq $(name)
migrate_up:
	migrate -path $(MIGRATION_DIR) -database $(MIGRATION_URL) -verbose up
migrate_down:
	migrate -path $(MIGRATION_DIR) -database $(MIGRATION_URL) -verbose down 1
migrate_status:
	migrate -path $(MIGRATION_DIR) -database $(MIGRATION_URL) version

.PHONY: test
test:
	 go test ./...

.PHONY: lint
lint:
	$(GOBIN)/golangci-lint/$(GOLANG_CI_VERSION)/golangci-lint run --config .golangci.yml

.PHONY: tools
tools:
	@if [ ! -f $(GOBIN)/golangci-lint ]; then\
		echo "Installing golangci-lint $(GOLANG_CI_VERSION)";\
		GOBIN=$(GOBIN)/golangci-lint/$(GOLANG_CI_VERSION) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANG_CI_VERSION);\
		echo "Done";\
	fi

	@if [ ! -f $(GOBIN)/gci ]; then\
		echo "Installing gci";\
		GOBIN=$(GOBIN)/gci/$(GCI_VERSION) go install github.com/daixiang0/gci@$(GCI_VERSION);\
		echo "Done";\
	fi

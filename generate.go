package generate

//go:generate $MOCKGEN_PATH -source=internal/domain/services/task.go -destination=internal/domain/services/task_mock.go --package=services

//go:generate $MOCKGEN_PATH -source=internal/infrastructure/storages/psql_store/task.go -destination=internal/infrastructure/storages/psql_store/task_mock.go --package=psql_store
//go:generate $MOCKGEN_PATH -source=internal/infrastructure/storages/psql_store/curency.go -destination=internal/infrastructure/storages/psql_store/curency_mock.go --package=psql_store

//go:generate $MOCKGEN_PATH -source=internal/infrastructure/storages/transaction/transaction_manager.go -destination=internal/infrastructure/storages/transaction/transaction_manager_mock.go --package=transaction

//go:generate $MOCKGEN_PATH -source=internal/events/event.go -destination=internal/events/event_mock.go --package=events

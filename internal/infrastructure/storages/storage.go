package storages

import "github.com/blr-coder/tasks-svc/internal/infrastructure/storages/psql_store"

type IStorage interface {
	psql_store.ITaskStorage
	//StatusHistoryStorage
	//WithTransaction
}

package psql_store

type IStorage interface {
	ITaskStorage
	//StatusHistoryStorage
	//WithTransaction
}

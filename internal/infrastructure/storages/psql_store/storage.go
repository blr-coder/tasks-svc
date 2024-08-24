package psql_store

type IStorage interface {
	ITaskStorage
	ICurrencyStorage
	//StatusHistoryStorage
	//WithTransaction
}

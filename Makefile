MIGRATION_DIR = internal/infrastructure/storages/psql_store/migrations
MIGRATION_URL = "postgres://task_svc_db_user:task_svc_db_user_pass@localhost:5232/task_svc_db?sslmode=disable"

migrate_new: # pass name as parameter, example - make migrate_new name=add_task_table
	migrate create -ext sql -dir $(MIGRATION_DIR)/ -seq $(name)
migrate_up:
	migrate -path $(MIGRATION_DIR) -database $(MIGRATION_URL) -verbose up
migrate_down:
	migrate -path $(MIGRATION_DIR) -database $(MIGRATION_URL) -verbose down 1

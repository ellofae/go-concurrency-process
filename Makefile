.PHONY: migrate_up
migrate_up:
	migrate -path ./migrations/ -database 'postgres://postgres:password@database:5432/process?sslmode=disable' up

.PHONY: migrate_down
migrate_down:
	migrate -path ./migrations/ -database 'postgres://postgres:password@database:5432/process?sslmode=disable' down

.PHONY: migrate_up
migrate_up:
	golang-migrate -path ./migrations/ -database 'postgres://postgres:password@localhost:5432/process?sslmode=disable' up

.PHONY: migrate_down
migrate_down:
	golang-migrate -path ./migrations/ -database 'postgres://postgres:password@localhost:5432/process?sslmode=disable' down

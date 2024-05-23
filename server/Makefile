docker-run-postgres:
	docker run --name postgres -e POSTGRES_PASSWORD=admin -e POSTGRES_USER=admin -p 5432:5432  -d postgres

migrate-up:
	migrate -path db/migration -database "postgresql://admin:admin@localhost:5432/chat-app?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgresql://admin:admin@localhost:5432/chat-app?sslmode=disable" -verbose down

.PHONY: docker-run-postgres migrateup migratedown
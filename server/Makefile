migrateup:
	migrate -path db/migration -database "postgresql://admin:admin@localhost:5432/chat-app?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://admin:admin@localhost:5432/chat-app?sslmode=disable" -verbose down

.PHONY: migrateup migratedown
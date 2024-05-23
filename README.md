# chat-app-golang
- Run postgres using Docker
```
make docker-run-postgres
```
- Init schema
```
make migrate-up
```
- Run server: in server folder
```
go run main.go
```
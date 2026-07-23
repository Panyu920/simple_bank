postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=panyu -e POSTGRES_PASSWORD=panyu -d 57c72fd2a128	

createdb:
	docker exec -it postgres12 createdb --username=panyu --owner=panyu simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://panyu:panyu@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://panyu:panyu@localhost:5432/simple_bank?sslmode=disable" -verbose down

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)
sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go simple_bank/db/sqlc Store

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl

.PHONY:postgres createdb dropdb migrateup migratedown sqlc server mock proto
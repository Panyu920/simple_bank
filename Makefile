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

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY:postgres createdb dropdb migrateup migratedown sqlc
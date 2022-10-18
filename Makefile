postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

restartpg:
	docker restart postgres12

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root school

dropdb:
	docker exec -it postgres12 dropdb school

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/school?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/school?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc
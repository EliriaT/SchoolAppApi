DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

network:

	docker network create school-network

postgres:
	docker run --name postgres12 --network school-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

restartpg:
	docker restart postgres12

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root school

dropdb:
	docker exec -it postgres12 dropdb school

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migrateversion:
	migrate -path db/migration/ -database "$(DB_URL)" force 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/EliriaT/SchoolAppApi/db/sqlc Store

deployeks:
	kubectl apply -f eks/deployment.yaml


.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migrateversion

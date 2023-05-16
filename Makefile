createdb:
	docker exec -it postgres-13 createdb  --username=skygazer --owner=skygazer progressme

migrateup:
	migrate -path db/sql/migration -database "postgresql://skygazer:hamdalah@localhost:5232/progressme?sslmode=disable" -verbose up

migratedown:
	migrate -path db/sql/migration -database "postgresql://skygazer:hamdalah@localhost:5232/progressme?sslmode=disable" -verbose down


makeFileDir := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

sqlc:
	docker run --rm -v $(makeFileDir):/src -w /src kjconroy/sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc
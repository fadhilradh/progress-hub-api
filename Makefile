createdb:
	docker exec -it postgres-13 createdb  --username=skygazer --owner=skygazer progressme

migrateup:
	migrate -path db/sql/migration -database "postgresql://skygazer:hamdalah@172.26.176.1:5232/progressme?sslmode=disable" -verbose up

migratedown:
	migrate -path db/sql/migration -database "postgresql://skygazer:hamdalah@172.26.176.1:5232/progressme?sslmode=disable" -verbose down


makeFileDir := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

sqlc:
	sudo docker run --rm -v $(makeFileDir):/src -w /src kjconroy/sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc
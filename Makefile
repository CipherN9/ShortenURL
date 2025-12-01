PG_ADDR=your-database

migrate-up:
	go run -tags 'postgres,file' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0 \
    		-path db/migrations -database "$(PG_ADDR)" up

migrate-down:
	go run -tags 'postgres,file' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0 \
		-path db/migrations -database "$(PG_ADDR)" down 1

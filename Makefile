up:
	docker compose up -d

down:
	docker compose down -v

logs:
	docker compose logs -f

test:
	DB_DSN="postgres://test:test@localhost:54321/notes_test?sslmode=disable" go test -v ./integration/...

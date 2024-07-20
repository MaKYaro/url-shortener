.PHONY: cover
cover:
	go test -short -count=1 -race -coverprofile=c.out ./...
	go tool cover -html=c.out -o coverage.html 
	firefox coverage.html

.PHONY: clean-cover
clean-cover:
	rm c.out coverage.html

.PHONY: start-postgres
start-postgres:
	sudo service postgresql start

.PHONY: stop-postgres
stop-postgres:
	sudo service postgresql stop

.PHONY: run-shortener
run-shortener:
	go run ./cmd/url-shortener/main.go --config-path=./config/local.json

.PHONY: up-migrations
up-migrations:
	goose -dir ./migrations postgres "postgresql://postgres:postgres@localhost:5432/url_shortener?sslmode=disable" up

.PHONY: down-migrations
down-migrations:
	goose -dir ./migrations postgres "postgresql://postgres:postgres@localhost:5432/url_shortener?sslmode=disable" down

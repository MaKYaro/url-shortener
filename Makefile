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

.PHONY: build-bin
build-bin:
	go build -o url-shortener-bin ./cmd/url-shortener/main.go

.PHONY: run-bin
run-bin:
	./url-shortener-bin --config-path=./config/local.json

.PHONY: docker-us-rebuild
docker-us-rebuild:
	sudo docker rm us-container
	sudo docker image rm us-image
	sudo docker build -t us-image .

.PHONY: docker-us-run
docker-us-run:
	sudo docker run --name=us-container -p 8082:8082 us-image

.PHONY: docker-us-start
docker-us-start:
	sudo docker start us-container 

.PHONY: docker-us-stop
docker-us-stop:
	sudo docker stop us-container 

.PHONY: docker-db-start
docker-db-start:
	sudo docker start us-postgres

.PHONY: docker-db-stop
docker-db-stop:
	sudo docker stop us-postgres
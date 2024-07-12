.PHONY: cover
cover:
	go test -short -count=1 -race -coverprofile=c.out ./...
	go tool cover -html=c.out -o coverage.html 
	firefox coverage.html

.PHONY: clean-cover
clean-cover:
	rm c.out coverage.html
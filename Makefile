build:
	@go build  -o bin/go-nba 
run: build
	@./bin/go-nba
test:
	@go test -v ./...


all:build

build:
	go build -v -o bin/dukes main.go
	go build -v -o bin/dukes-cli cmd/cli/main.go
	
docker:
	docker build --rm -t dukes:0.1.0 .
	
test:
	@go test -v ./...

run:
	@go run main.go

run-cli:
	@go run cmd/cli/main.go

clean:
	$(RM) -r bin
	$(RM) -r raft
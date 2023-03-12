
all:build

build:
	go build -o bin/dukes main.go
	go build -o bin/dukes-cli cmd/cli/main.go
	
test:
	@go test -v ./...

run:
	@go run main.go

run-cli:
	@go run cmd/cli/main.go

clean:
	$(RM) bin/*
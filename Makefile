# Copyright (c) 2023 Robeto Ughi
# 
# This software is released under the MIT License.
# https://opensource.org/licenses/MIT


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
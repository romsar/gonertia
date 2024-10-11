test:
	go test -race -count 1 ./...

supertest:
	go test -race -count 100 ./...

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run ./... --fix
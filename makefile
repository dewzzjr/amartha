mod:
	@go mod tidy
	@go mod vendor

test:
	@go vet ./...
	@go test -timeout 30s -cover -coverpkg=all ./...

build:
	@go build -o bin/reconcile .

run:
	@bin/reconcile \
		-f test/data/system.csv \
		-b test/data/bca.csv \
		-b test/data/bni.csv \
		-v

.PHONY: test
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
		--start 2024-06-20 \
		--end 2024-06-21

debug:
	@bin/reconcile -v \
		-f test/data/system.csv \
		-b test/data/bca.csv \
		-b test/data/bni.csv \
		--start 2024-06-20 \
		--end 2024-06-21

.PHONY: test
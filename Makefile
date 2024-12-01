test:  
	go clean -testcache
	go test -timeout 1m -cover ./...

lint: 
	golangci-lint run --allow-parallel-runners -c ./.golangci-lint.yaml --fix ./...

.PHONY: test lint
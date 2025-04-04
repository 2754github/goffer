init:
	go install golang.org/x/tools/cmd/godoc@latest
	go install golang.org/x/tools/cmd/goimports@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $$(go env GOPATH)/bin v2.1.2

doc:
	open http://localhost:6060/pkg/github.com/2754github/goffer/
	godoc -http=:6060

tidy:
	go mod tidy

format: tidy
	go fmt ./...
	goimports -local github.com/2754github/goffer -w .

lint: format
	go vet ./...
	$$(go env GOPATH)/bin/golangci-lint run

test: lint
	go clean -testcache
	go test -v ./...

.PHONY: init doc tidy format lint test

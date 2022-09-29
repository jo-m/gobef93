.PHONY: test format lint check

format:
	gofmt -w .
	go mod tidy

test:
	go test -v -race ./...

lint:
	gofmt -l .; test -z "$$(gofmt -l .)"

	go vet ./...
	
	# go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck -checks=all ./...
	
	# go install github.com/mgechev/revive@latest
	revive ./...

	# go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec scan -checks=all ./...

check: lint test

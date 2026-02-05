.PHONY: test run clean

test:
	go test ./... -v

run:
	go run ./cmd/textpad

clean:
	go clean -testcache

build: fmt test lint
	mkdir -p bin && go build -o ./bin/rgt ./cmd/rgt

test:
	go test ./... -skip 'TestFast|TestSlow'

test-slow: build
	./bin/rgt start --test-name TestSlow

test-fast: build
	./bin/rgt start --test-name TestFast

fmt:
	go fmt ./...

lint:
	golint ./...

run: build
	./bin/rgt start

install-local-dev: build
	sudo cp ./bin/rgt /usr/local/bin


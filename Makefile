build: fmt test lint
	mkdir -p bin && go build -o ./bin/rgt ./cmd/rgt

test:
	go test ./...

fmt:
	go fmt ./...

lint:
	golint ./...

run: build
	./bin/rgt start

install-local-dev: build
	sudo cp ./bin/rgt /usr/local/bin


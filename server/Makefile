.PHONY: build run dev clean migrate test deploy

BINARY=jianghu-server
BUILD_DIR=bin

build:
	CGO_ENABLED=1 go build -o $(BUILD_DIR)/$(BINARY) ./cmd/server

run: build
	./$(BUILD_DIR)/$(BINARY)

dev:
	@which air > /dev/null 2>&1 || go install github.com/air-verse/air@latest
	air -c .air.toml

migrate:
	CGO_ENABLED=1 go run ./cmd/server -migrate

test:
	CGO_ENABLED=1 go test -v ./...

clean:
	rm -rf $(BUILD_DIR)
	rm -f data/*.db

lint:
	golangci-lint run ./...

docker-build:
	docker build -t jianghu-server .

docker-run:
	docker run -p 8080:8080 -v $(PWD)/data:/app/data jianghu-server

deploy: build
	@echo "Deploying $(BINARY)..."
	scp $(BUILD_DIR)/$(BINARY) deploy@server:/opt/jianghu-server/
	scp config.yaml deploy@server:/opt/jianghu-server/
	ssh deploy@server "sudo systemctl restart jianghu-server"

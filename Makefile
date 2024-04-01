# Commands
.PHONY: env
env:
	@$(eval SHELL:=/bin/bash)
	@cp .env.sample .env
	@echo "SECRET_KEY=$$(openssl rand -hex 32)" >> .env

.PHONY: build
build:
	go build -o ./cmd/auth

.PHONY: run
run:
	go run ./cmd/auth

.PHONY: compose-up
compose-up:
	docker-compose -f docker-compose.yml up -d --remove-orphans

.PHONY: compose-down
compose-down:
	docker-compose down --remove-orphans

.PHONY: unit-test
unit-test:
	go test ./...

.PHONY: cover
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	rm coverage.out

.PHONY: linter
linter:
	golangci-lint run

.PHONY: swagger
swagger:
	swag init -g internal/app/app.go --parseInternal --parseDependency

.PHONY: proto
proto:
	rm -f pkg/pb/*.go
	protoc -I proto proto/*.proto --go_out=./pkg/pb --go_opt=paths=source_relative \
	--go-grpc_out=./pkg/pb --go-grpc_opt=paths=source_relative \
#	rm -f doc/swagger/*.swagger.json
#	--grpc-gateway_out=pkg/pb --grpc-gateway_opt=paths=source_relative \
#	--openapiv2_out=doc/swagger \
#	--openapiv2_opt logtostderr=true --openapiv2_opt use_go_templates=true \
#	proto/shortener.proto
generate-proto:
	@protoc -I ./api -I ./api/google \
	--go_out=./pkg/pb --go_opt=paths=source_relative \
	--go-grpc_out=./pkg/pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=./pkg/pb --grpc-gateway_opt=paths=source_relative \
	./api/movie.proto

lint:
	@golangci-lint run

build:
	@go build -o bin/app cmd/app/main.go

run:
	@go run cmd/app/main.go

compose-up:
	docker compose up -d

compose-down:
	docker compose down --remove-orphans
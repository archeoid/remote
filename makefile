server:
	@go build -o bin cmd/server/server.go
	@go build -o bin cmd/client/client.go

flatbuffer: api/message.fbs
	@echo "Compiling flatbuffer schema."
	@flatc --go api/message.fbs ''
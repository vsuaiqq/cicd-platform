.PHONY: proto clean auth-proto projects-proto ai-proto

proto: auth-proto projects-proto ai-proto
	@echo "All proto files generated successfully!"

auth-proto:
	@echo "Generating gRPC code for auth domain..."
	@mkdir -p gen/auth
	@protoc --proto_path=api/auth --go_out=gen/auth --go_opt=paths=source_relative \
		--go-grpc_out=gen/auth --go-grpc_opt=paths=source_relative \
		auth.proto
	@echo "Auth proto files generated in gen/auth/"

projects-proto:
	@echo "Generating gRPC code for projects domain..."
	@mkdir -p gen/projects
	@protoc --proto_path=api/projects --go_out=gen/projects --go_opt=paths=source_relative \
		--go-grpc_out=gen/projects --go-grpc_opt=paths=source_relative \
		projects.proto
	@echo "Projects proto files generated in gen/projects/"

ai-proto:
	@echo "Generating gRPC code for ai domain..."
	@mkdir -p gen/ai
	@protoc --proto_path=api/ai --go_out=gen/ai --go_opt=paths=source_relative \
		--go-grpc_out=gen/ai --go-grpc_opt=paths=source_relative \
		ai.proto
	@echo "AI proto files generated in gen/ai/"

clean:
	@rm -rf gen/
	@echo "Cleaned generated proto files"

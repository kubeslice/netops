.PHONY: compile
compile: ## Compile the proto file.
	protoc -I pkg/proto/ pkg/proto/netop.proto --go_out=paths=source_relative:pkg/proto --go-grpc_out=pkg/proto --go-grpc_opt=paths=source_relative

.PHONY: kubeslice-netops
kubeslice-netops: ## Build and run kubeslice-netops.
	go build -race -ldflags "-s -w" -o bin/kubeslice-netops main.go
	bin/kubeslice-netops

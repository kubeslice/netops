# VERSION defines the project version for the bundle.
# Update this value when you upgrade the version of your project.
VERSION ?= latest-stable

IMG ?= docker.io/aveshasystems/netops:$(VERSION)

.PHONY: compile
compile: ## Compile the proto file.
	protoc -I pkg/proto/ pkg/proto/netop.proto --go_out=paths=source_relative:pkg/proto --go-grpc_out=pkg/proto --go-grpc_opt=paths=source_relative

.PHONY: kubeslice-netops
kubeslice-netops: ## Build and run kubeslice-netops.
	go build -race -ldflags "-s -w" -o bin/kubeslice-netops main.go

.PHONY: docker-build
docker-build: ## Build docker image with the manager.
	docker buildx create --name container --driver=docker-container || true
	docker build --builder container --platform linux/amd64,linux/arm64 -t ${IMG} .

.PHONY: docker-push
docker-push: ## Push docker image with the manager.
	docker buildx create --name container --driver=docker-container || true
	docker build --push --builder container --platform linux/amd64,linux/arm64 -t ${IMG} .

.PHONY: chart-deploy
chart-deploy:
	## Deploy the artifacts using helm
	## Usage: make chart-deploy VALUESFILE=[valuesfilename]
	helm upgrade --install kubeslice-worker -n kubeslice-system avesha/kubeslice-worker -f ${VALUESFILE}

# VERSION defines the project version for the bundle.
# Update this value when you upgrade the version of your project.
VERSION ?= latest-stable

.PHONY: compile
compile: ## Compile the proto file.
	protoc -I pkg/proto/ pkg/proto/netop.proto --go_out=paths=source_relative:pkg/proto --go-grpc_out=pkg/proto --go-grpc_opt=paths=source_relative

.PHONY: kubeslice-netops
kubeslice-netops: ## Build and run kubeslice-netops.
	go build -race -ldflags "-s -w" -o bin/kubeslice-netops main.go

.PHONY: docker-build
docker-build: kubeslice-netops
	docker build -t kubeslice-netop:${VERSION} --build-arg PLATFORM=amd64 . && docker tag kubeslice-netop:${VERSION} docker.io/rahulsawra/netops:${VERSION}

.PHONY: chart-deploy
chart-deploy:
	## Deploy the artifacts using helm
	## Usage: make chart-deploy VALUESFILE=[valuesfilename]
	helm upgrade --install kubeslice-worker -n kubeslice-system avesha/kubeslice-worker -f ${VALUESFILE}

.PHONY: docker-push
docker-push:
	docker push docker.io/rahulsawra/netops:${VERSION}
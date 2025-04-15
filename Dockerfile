##########################################################
#Dockerfile
#Copyright (c) 2022 Avesha, Inc. All rights reserved.
#
#SPDX-License-Identifier: Apache-2.0
#
#Licensed under the Apache License, Version 2.0 (the "License");
#you may not use this file except in compliance with the License.
#You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS,
#WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#See the License for the specific language governing permissions and
#limitations under the License.
##########################################################

FROM golang:1.24.2 AS gobuilder

ARG TARGETPLATFORM
ARG TARGETARCH
ARG TARGETOS

# Set the Go source path
WORKDIR /kubeslice/kubeslice-netops/
COPY go.mod go.sum ./
ADD vendor vendor
COPY . .
# Build the binary.esah
RUN go env -w GOPRIVATE=github.com/kubeslice && \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} GO111MODULE=on go build -mod=vendor -a -o bin/kubeslice-netops main.go

# Build reduced image from base alpine
FROM alpine:3.21

# Add the necessary pakages:
# tc - is needed for traffic control and shaping on the kubeslice-netops.  it is part of the iproute2
RUN apk add --no-cache ca-certificates \
    iproute2
# Run the kubeslice-netops binary.
WORKDIR /kubeslice

# Copy our static executable.
COPY --from=gobuilder /kubeslice/kubeslice-netops/bin/kubeslice-netops .
EXPOSE 5000
EXPOSE 8080
# Or could be CMD
ENTRYPOINT ["./kubeslice-netops"]

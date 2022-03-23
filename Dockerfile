######################################################################
# Project: kubeslice-netops
# File: kubeslice-netops/Dockerfile
# Created: 02/15/2022
#
# Avesha LLC
# Copyright (c) Avesha LLC. 2019, 2020, 2021, 2022
######################################################################
ARG PLATFORM
FROM ${PLATFORM}/golang:1.17.7-alpine3.15 as gobuilder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git make build-base

# Set the Go source path
WORKDIR /kubeslice/kubeslice-netops/
COPY . .
# Build the binary.esah
RUN go mod download && \
    go env -w GOPRIVATE=bitbucket.org/realtimeai && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/kubeslice-netops main.go

# Build reduced image from base alpine
FROM ${PLATFORM}/alpine:3.15

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

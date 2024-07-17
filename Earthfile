VERSION 0.6
# renovate: datasource=docker depName=quay.io/kairos/osbuilder-tools versioning=semver-coerced
ARG OSBUILDER_VERSION=v0.200.12
ARG OSBUILDER_IMAGE=quay.io/kairos/osbuilder-tools:$OSBUILDER_VERSION
# renovate: datasource=docker depName=golangci/golangci-lint
ARG GOLINT_VERSION=v1.57.2
# renovate: datasource=docker depName=golang
ARG GO_VERSION=1.20-bookworm
ARG GOLANG_VERSION=1.22

version:
    FROM +go-deps
    COPY . ./
    RUN --no-cache echo $(git describe --always --tags --dirty) > VERSION
    RUN --no-cache echo $(git describe --always --dirty) > COMMIT
    ARG VERSION=$(cat VERSION)
    ARG COMMIT=$(cat COMMIT)
    SAVE ARTIFACT VERSION VERSION
    SAVE ARTIFACT COMMIT COMMIT

go-deps:
    ARG GO_VERSION
    FROM gcr.io/spectro-images-public/golang:${GOLANG_VERSION}-alpine
    # RUN apt-get update && apt-get install -y rsync gcc bash git
    WORKDIR /build
    COPY . .
    RUN go mod tidy --compat=1.22
    RUN go mod download
    RUN go mod verify

test:
    FROM +go-deps
    RUN go run github.com/onsi/ginkgo/v2/ginkgo --race --covermode=atomic --coverprofile=coverage.out -p -r ./...
    SAVE ARTIFACT coverage.out AS LOCAL coverage.out

golint:
    ARG GOLINT_VERSION
    FROM golangci/golangci-lint:$GOLINT_VERSION
    WORKDIR /build
    COPY . .
    RUN go mod tidy --compat=1.19
    RUN golangci-lint run -v

build-immucore:
    FROM +go-deps
    COPY +version/VERSION ./
    COPY +version/COMMIT ./
    ARG VERSION=$(cat VERSION)
    ARG COMMIT=$(cat COMMIT)
    # ARG LDFLAGS="-s -w -X github.com/kairos-io/immucore/internal/version.version=$VERSION -X github.com/kairos-io/immucore/internal/version.gitCommit=$COMMIT"
    # RUN echo ${LDFLAGS}
    # RUN CGO_ENABLED=0 go build -o immucore -ldflags "${LDFLAGS}"
    ARG LDFLAGS="-w -X github.com/kairos-io/immucore/internal/version.version=v0.1.34_spectro"
    RUN go-build-fips.sh -a -o immucore
    SAVE ARTIFACT immucore immucore AS LOCAL build/immucore-$VERSION

# Alias for ease of use
build:
    BUILD +build-immucore

framework:
    FROM quay.io/kairos/framework:v2.7.41-fips
    
    COPY +build-immucore/immucore* /usr/bin/immucore
    
    SAVE IMAGE --push gcr.io/spectro-dev-public/nianyu/framework:v2.7.41-fips-spectro
    
dracut-artifacts:
    FROM scratch
    WORKDIR /build
    COPY --dir dracut/28immucore .
    COPY dracut/10-immucore.conf .
    SAVE ARTIFACT 28immucore 28immucore
    SAVE ARTIFACT 10-immucore.conf 10-immucore.conf
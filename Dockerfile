ARG GO_VERSION=latest

ARG DISTROLESS_IMAGE=gcr.io/distroless/static
ARG DISTROLESS_IMAGE_TAG=nonroot

ARG NGT_BUILD_OPTIONS="-DNGT_AVX_DISABLED=ON"
ARG VERSION="unknown"

FROM golang:${GO_VERSION} AS builder
ARG NGT_BUILD_OPTIONS
ARG VERSION

ENV ORG rinx
ENV REPO alvd

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    build-essential \
    cmake \
    curl \
    unzip \
    git \
    gcc \
    g++ \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}
COPY go.mod .
COPY go.sum .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/cmd
COPY cmd .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}/pkg/alvd
COPY pkg/alvd .

WORKDIR ${GOPATH}/src/github.com/${ORG}/${REPO}
COPY Makefile .

RUN make NGT_BUILD_OPTIONS="${NGT_BUILD_OPTIONS}" ngt/install
RUN make VERSION="${VERSION}" cmd/alvd/alvd \
    && cp cmd/alvd/alvd /alvd

FROM ${DISTROLESS_IMAGE}:${DISTROLESS_IMAGE_TAG}
LABEL maintainer "rintaro okamura <rintaro.okamura@gmail.com>"

COPY --from=builder /alvd /alvd

USER nonroot:nonroot

ENTRYPOINT ["/alvd"]

FROM mcr.microsoft.com/vscode/devcontainers/go:1 AS base

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    cmake \
    curl \
    g++ \
    gawk \
    gcc \
    git \
    jq \
    sed \
    zip \
    unzip \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR ${GOPATH}

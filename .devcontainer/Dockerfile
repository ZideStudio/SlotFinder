FROM mcr.microsoft.com/devcontainers/base:ubuntu

ARG TARGETARCH
ENV GOTOOLCHAIN=local

# Go
RUN curl -fsSL https://go.dev/dl/go1.25.5.linux-${TARGETARCH}.tar.gz \
    | tar -C /usr/local -xz

# Configure Go environment
ENV GOPATH=/go
ENV GOBIN=/go/bin
ENV PATH="/usr/local/go/bin:/go/bin:${PATH}"

# Node
RUN curl -fsSL https://nodejs.org/dist/v24.13.0/node-v24.13.0-linux-${TARGETARCH}.tar.xz \
    | tar -xJ -C /usr/local --strip-components=1

# Utils
RUN apt-get update && apt-get install -y \
    git \
    make \
    curl \
    unzip \
    postgresql-client \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Create Go directories
RUN mkdir -p /go/bin /go/pkg /go/src

# Go tools (LSP, formatter, debugger)
RUN go install golang.org/x/tools/gopls@v0.21.0 && \
    go install github.com/go-delve/delve/cmd/dlv@v1.26.0 && \
    go install golang.org/x/tools/cmd/goimports@latest && \
    go install github.com/air-verse/air@latest

# TypeScript/JavaScript LSP
RUN npm install -g typescript typescript-language-server vscode-langservers-extracted

# Create prettier wrapper script for Zed
RUN echo '#!/bin/sh\ncd /workspace/front\nexec npx prettier "$@"' > /usr/local/bin/prettier-wrapper && \
    chmod +x /usr/local/bin/prettier-wrapper

# Note: prettier and oxlint will be used from the front project's node_modules
# This ensures version consistency with the project

WORKDIR /workspace

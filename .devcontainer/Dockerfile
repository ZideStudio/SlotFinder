FROM mcr.microsoft.com/devcontainers/base:ubuntu

ARG TARGETARCH
ENV GOTOOLCHAIN=local

# Go
RUN curl -fsSL https://go.dev/dl/go1.25.5.linux-${TARGETARCH}.tar.gz \
    | tar -C /usr/local -xz

ENV PATH="/usr/local/go/bin:${PATH}"

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

# Go tools
RUN go install golang.org/x/tools/gopls@v0.21.0 \
    && go install github.com/go-delve/delve/cmd/dlv@v1.26.0

# Front tools (LSP side)
RUN npm install -g prettier oxlint

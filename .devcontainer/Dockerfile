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
    git \
    && rm -rf /var/lib/apt/lists/*

# Install lazygit
ARG LAZYGIT_VERSION=0.44.1
RUN curl -Lo lazygit.tar.gz "https://github.com/jesseduffield/lazygit/releases/download/v${LAZYGIT_VERSION}/lazygit_${LAZYGIT_VERSION}_Linux_${TARGETARCH}.tar.gz" && \
    tar xf lazygit.tar.gz lazygit && \
    install lazygit /usr/local/bin && \
    rm lazygit.tar.gz lazygit

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

# Create shell aliases
RUN echo 'alias lg="lazygit"' >> /home/vscode/.bashrc && \
    echo 'alias lg="lazygit"' >> /home/vscode/.zshrc 2>/dev/null || true

# Create Git initialization script
RUN echo '#!/bin/sh\n\
    git config --global --add safe.directory /workspace 2>/dev/null || true\n\
    git config --global --add safe.directory "*" 2>/dev/null || true\n\
    ' > /usr/local/bin/init-git.sh && \
    chmod +x /usr/local/bin/init-git.sh

# Run git initialization at container start
RUN echo '#!/bin/sh\n\
    /usr/local/bin/init-git.sh\n\
    exec "$@"\n\
    ' > /usr/local/bin/docker-entrypoint.sh && \
    chmod +x /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["sleep", "infinity"]

WORKDIR /workspace

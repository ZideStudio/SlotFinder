FROM mcr.microsoft.com/devcontainers/base:ubuntu

ARG TARGETARCH
ENV GOTOOLCHAIN=local
ENV GOPATH=/go
ENV GOBIN=/go/bin
ENV PATH="/usr/local/go/bin:/go/bin:${PATH}"

# System
USER root

RUN apt-get update && apt-get install -y \
    git make curl unzip postgresql-client build-essential \
    && rm -rf /var/lib/apt/lists/*

RUN if [ "$TARGETARCH" = "amd64" ]; then \
    echo "x64" > /tmp/node_arch && \
    echo "x86_64" > /tmp/lg_arch; \
    else \
    echo "$TARGETARCH" > /tmp/node_arch && \
    echo "$TARGETARCH" > /tmp/lg_arch; \
    fi

# Go
COPY back/go.mod /tmp/go.mod
RUN GO_VERSION=$(grep -oP '^go\s+\K[\d.]+' /tmp/go.mod) && \
    curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${TARGETARCH}.tar.gz" | tar -C /usr/local -xz

# Git
ARG LAZYGIT_VERSION=0.44.1
RUN if [ "$TARGETARCH" = "amd64" ]; then LG_ARCH="x86_64"; else LG_ARCH=$TARGETARCH; fi && \
    curl -Lo lazygit.tar.gz "https://github.com/jesseduffield/lazygit/releases/download/v${LAZYGIT_VERSION}/lazygit_${LAZYGIT_VERSION}_Linux_${LG_ARCH}.tar.gz" && \
    tar xf lazygit.tar.gz lazygit && \
    install lazygit /usr/local/bin && \
    rm lazygit.tar.gz lazygit

RUN git config --system --add safe.directory /workspace

# Scripts
RUN mkdir -p /go/bin /go/pkg /go/src && chown -R vscode:vscode /go

RUN echo '#!/bin/sh\nexec "$@"' > /usr/local/bin/docker-entrypoint.sh && \
    chmod +x /usr/local/bin/docker-entrypoint.sh

RUN echo '#!/bin/bash\n\
    export VOLTA_HOME="/home/vscode/.volta"\n\
    export PATH="$VOLTA_HOME/bin:$PATH"\n\
    cd /workspace/front 2>/dev/null || cd .\n\
    if [ -f "./node_modules/.bin/prettier" ]; then\n\
    exec ./node_modules/.bin/prettier "$@"\n\
    else\n\
    exec npx prettier "$@"\n\
    fi' > /usr/local/bin/prettier && \
    chmod +x /usr/local/bin/prettier

# Tools
USER vscode

# Volta
ENV VOLTA_HOME="/home/vscode/.volta"
ENV PATH="${VOLTA_HOME}/bin:${PATH}"
RUN curl https://get.volta.sh | bash -s -- --skip-setup

COPY --chown=vscode:vscode front/package.json /tmp/package.json
RUN NODE_VERSION=$(grep -oP '"node":\s*"\K[^"]+' /tmp/package.json | head -1) && \
    volta install node@${NODE_VERSION:-latest} && \
    volta install typescript typescript-language-server vscode-langservers-extracted

# Go
RUN go install golang.org/x/tools/gopls@v0.21.0 && \
    go install github.com/go-delve/delve/cmd/dlv@v1.26.0 && \
    go install golang.org/x/tools/cmd/goimports@latest && \
    go install github.com/air-verse/air@latest

RUN echo 'alias lg="lazygit"' >> /home/vscode/.bashrc

WORKDIR /workspace
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["sleep", "infinity"]

FROM mcr.microsoft.com/devcontainers/base:ubuntu

ARG TARGETARCH
ENV GOTOOLCHAIN=local

RUN if [ "$TARGETARCH" = "amd64" ]; then \
    echo "x64" > /tmp/arch_node && \
    echo "x86_64" > /tmp/arch_lazygit; \
    else \
    echo "$TARGETARCH" > /tmp/arch_node && \
    echo "$TARGETARCH" > /tmp/arch_lazygit; \
    fi

# Go
RUN curl -fsSL https://go.dev/dl/go1.25.5.linux-${TARGETARCH}.tar.gz \
    | tar -C /usr/local -xz

ENV GOPATH=/go
ENV GOBIN=/go/bin
ENV PATH="/usr/local/go/bin:/go/bin:${PATH}"

# Volta
ENV VOLTA_HOME="/home/vscode/.volta"
ENV PATH="${VOLTA_HOME}/bin:${PATH}"
USER vscode

RUN curl https://get.volta.sh | bash -s -- --skip-setup

COPY --chown=vscode:vscode front/package.json /tmp/package.json

RUN NODE_VERSION=$(grep -oP '"node":\s*"\K[^"]+' /tmp/package.json | head -1) && \
    echo "Installing Node ${NODE_VERSION} from package.json..." && \
    volta install node@${NODE_VERSION} && \
    volta install typescript typescript-language-server vscode-langservers-extracted


# Use root for system tools
USER root

# Utils
RUN apt-get update && apt-get install -y \
    git \
    make \
    curl \
    unzip \
    postgresql-client \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Lazygit
ARG LAZYGIT_VERSION=0.44.1
RUN LG_ARCH=$(cat /tmp/arch_lazygit) && \
    curl -Lo lazygit.tar.gz "https://github.com/jesseduffield/lazygit/releases/download/v${LAZYGIT_VERSION}/lazygit_${LAZYGIT_VERSION}_Linux_${LG_ARCH}.tar.gz" && \
    tar xf lazygit.tar.gz lazygit && \
    install lazygit /usr/local/bin && \
    rm lazygit.tar.gz lazygit

# Permissions & Go Tools
RUN mkdir -p /go/bin /go/pkg /go/src && \
    chown -R vscode:vscode /go

RUN go install golang.org/x/tools/gopls@v0.21.0 && \
    go install github.com/go-delve/delve/cmd/dlv@v1.26.0 && \
    go install golang.org/x/tools/cmd/goimports@latest && \
    go install github.com/air-verse/air@latest

# Wrapper Prettier optimisé pour l'environnement Volta
RUN echo '#!/bin/bash\n\
    export VOLTA_HOME="/home/vscode/.volta"\n\
    export PATH="$VOLTA_HOME/bin:$PATH"\n\
    cd /workspace/front\n\
    if [ -f "./node_modules/.bin/prettier" ]; then\n\
    exec ./node_modules/.bin/prettier "$@"\n\
    else\n\
    exec prettier "$@"\n\
    fi' > /usr/local/bin/prettier && \
    chmod +x /usr/local/bin/prettier

# Aliases & Scripts
RUN echo 'alias lg="lazygit"' >> /home/vscode/.bashrc && \
    echo 'alias lg="lazygit"' >> /home/vscode/.zshrc 2>/dev/null || true

RUN echo '#!/bin/sh\n\
    git config --global --add safe.directory /workspace 2>/dev/null || true\n\
    git config --global --add safe.directory "*" 2>/dev/null || true\n\
    ' > /usr/local/bin/init-git.sh && \
    chmod +x /usr/local/bin/init-git.sh

RUN echo '#!/bin/sh\n\
    /usr/local/bin/init-git.sh\n\
    exec "$@"\n\
    ' > /usr/local/bin/docker-entrypoint.sh && \
    chmod +x /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["sleep", "infinity"]

WORKDIR /workspace

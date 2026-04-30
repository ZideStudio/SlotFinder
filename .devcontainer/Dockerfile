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
    xvfb sudo \
    && echo "vscode ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/vscode \
    && chmod 0440 /etc/sudoers.d/vscode \
    && mkdir -p /tmp/.X11-unix \
    && chmod 1777 /tmp/.X11-unix \
    && rm -rf /var/lib/apt/lists/*

RUN if [ "$TARGETARCH" = "amd64" ]; then \
    echo "x64" > /tmp/node_arch && \
    echo "x86_64" > /tmp/lg_arch; \
    else \
    echo "$TARGETARCH" > /tmp/node_arch && \
    echo "$TARGETARCH" > /tmp/lg_arch; \
    fi

COPY .devcontainer/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

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
    cd /workspace/front 2>/dev/null || cd .\n\
    if [ -f "./node_modules/.bin/prettier" ]; then\n\
    exec ./node_modules/.bin/prettier "$@"\n\
    else\n\
    exec npx prettier "$@"\n\
    fi' > /usr/local/bin/prettier && \
    chmod +x /usr/local/bin/prettier

# Tools
USER vscode

# Mise
ENV MISE_HOME="/home/vscode/.local/share/mise"
ENV MISE_TRUSTED_CONFIG_PATHS="/workspace"
ENV PATH="/home/vscode/.local/share/mise/shims:/home/vscode/.local/bin:/usr/local/go/bin:/go/bin:${PATH}"
RUN curl https://mise.jdx.dev/install.sh | sh

RUN mkdir -p /tmp/front-mise /tmp/back-mise

COPY --chown=vscode:vscode front/.mise.toml /tmp/front-mise/.mise.toml
COPY --chown=vscode:vscode back/.mise.toml /tmp/back-mise/.mise.toml

RUN cd /tmp/front-mise && /home/vscode/.local/bin/mise trust .mise.toml && /home/vscode/.local/bin/mise install
RUN cd /tmp/back-mise && /home/vscode/.local/bin/mise trust .mise.toml && /home/vscode/.local/bin/mise install

RUN /home/vscode/.local/bin/mise exec -C /tmp/front-mise -- npx playwright install --with-deps

# Go
RUN go install golang.org/x/tools/gopls@v0.21.1 && \
    go install github.com/go-delve/delve/cmd/dlv@v1.26.0 && \
    go install golang.org/x/tools/cmd/goimports@v0.44.0

RUN cat <<EOF >> /home/vscode/.bashrc
alias lg="lazygit"
export DISPLAY=:99
eval "\$(\$HOME/.local/bin/mise activate bash)"
EOF

WORKDIR /workspace
USER vscode
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["sleep", "infinity"]

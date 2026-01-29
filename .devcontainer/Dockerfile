FROM mcr.microsoft.com/devcontainers/go:1.24-bookworm

# Install node environment
# Remove yarn repo to avoid GPG key issues, then install packages
RUN rm -f /etc/apt/sources.list.d/yarn.list && \
    apt-get update && \
    apt-get install -y --no-install-recommends \
    nodejs graphviz npm \
    && rm -rf /var/lib/apt/lists/*

USER vscode

# Install go swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest
# Use standard golang build
FROM golang

# Clean all directories
RUN rm -rf "$GOPATH/src" && rm -rf "$GOPATH/pkg" && rm -rf "$GOPATH/bin"

# Copy our app files
ADD . "$GOPATH/src/github.com/OrfeasZ/ecr-token-refresh"

# Build and install our app
RUN go install github.com/OrfeasZ/ecr-token-refresh

# Expose the http health check server
EXPOSE 3277

# Set the config volume
VOLUME /opt/config/ecr-token-refresh

# Set our app entrypoint
ENTRYPOINT "$GOPATH/bin/ecr-token-refresh"
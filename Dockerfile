# Use standard golang build
FROM golang:1.8.3-alpine
# Copy our app files
ADD . "$GOPATH/src/ecr-token-refresh"

# Build and install our app
RUN go install ecr-token-refresh

# Secondary stage, binaries only
FROM alpine:latest  

RUN apk add -U ca-certificates

WORKDIR /bin/

# Copy the application binary
COPY --from=0 /go/bin/ecr-token-refresh .

# Expose the http health check server
EXPOSE 3277

# Set the config volume
VOLUME /opt/config/ecr-token-refresh

# Set our app entrypoint
ENTRYPOINT ["/bin/ecr-token-refresh"]

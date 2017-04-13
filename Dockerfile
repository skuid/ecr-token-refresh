FROM alpine

RUN apk -U add ca-certificates

COPY ecr-token-refresh /bin/ecr-token-refresh

VOLUME /opt/config/ecr-token-refresh

ENTRYPOINT ["/bin/ecr-token-refresh"]
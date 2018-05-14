FROM golang:1.10-alpine AS lsp-adapter
WORKDIR /go/src/github.com/sourcegraph/lsp-adapter
COPY . .
RUN CGO_ENABLED=0 GOBIN=/usr/local/bin go install github.com/sourcegraph/lsp-adapter

# The only tag available is latest.
# https://github.com/rcjsuen/dockerfile-language-server-nodejs/issues/222
# hadolint ignore=DL3007
FROM rcjsuen/docker-langserver:latest

RUN apk add --no-cache tini=0.16.1-r0
ENTRYPOINT ["tini", "--"]

COPY --from=lsp-adapter /usr/local/bin/lsp-adapter /usr/local/bin/
EXPOSE 8080
CMD ["lsp-adapter", "-proxyAddress=0.0.0.0:8080", "-glob=Dockerfile*", "/docker-langserver/bin/docker-langserver", "--stdio"]

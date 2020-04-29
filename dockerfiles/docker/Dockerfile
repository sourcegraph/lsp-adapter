FROM golang:1.14-alpine AS lsp-adapter
WORKDIR /go/src/github.com/sourcegraph/lsp-adapter
COPY . .
RUN CGO_ENABLED=0 GOBIN=/usr/local/bin go install github.com/sourcegraph/lsp-adapter

FROM rcjsuen/docker-langserver:0.0.19

RUN apk add --no-cache tini
ENTRYPOINT ["tini", "--"]

COPY --from=lsp-adapter /usr/local/bin/lsp-adapter /usr/local/bin/
EXPOSE 8080
CMD ["lsp-adapter", "-proxyAddress=0.0.0.0:8080", "-glob=Dockerfile*", "/docker-langserver/bin/docker-langserver", "--stdio"]

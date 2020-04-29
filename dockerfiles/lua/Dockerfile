FROM golang:1.14-alpine AS lsp-adapter
WORKDIR /go/src/github.com/sourcegraph/lsp-adapter
COPY . .
RUN CGO_ENABLED=0 GOBIN=/usr/local/bin go install github.com/sourcegraph/lsp-adapter

# 👀 Add steps here to build the language server itself 👀
# CMD ["echo", "🚨 This statement should be removed once you have added the logic to start up the language server! 🚨 Exiting..."]

FROM abaez/luarocks:lua5.1

RUN apk add --no-cache ca-certificates git tini

ENTRYPOINT ["/sbin/tini", "--"]

RUN luarocks install --server=http://luarocks.org/dev lua-lsp

COPY --from=lsp-adapter /usr/local/bin/lsp-adapter /usr/local/bin
EXPOSE 8080
# Modify this command to connect to the language server
CMD ["lsp-adapter", "--trace", "--proxyAddress=0.0.0.0:8080", "lua-lsp"]

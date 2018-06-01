FROM golang:1.10-alpine AS lsp-adapter
WORKDIR /go/src/github.com/sourcegraph/lsp-adapter
COPY . .
RUN CGO_ENABLED=0 GOBIN=/usr/local/bin go install github.com/sourcegraph/lsp-adapter

FROM clojure:lein-2.8.1-alpine
RUN apk add --no-cache ca-certificates=20171114-r0 git=2.15.2-r0 tini=0.16.1-r0

WORKDIR /tmp
RUN git clone https://github.com/snoe/clojure-lsp.git
WORKDIR /tmp/clojure-lsp
# https://github.com/snoe/clojure-lsp/commit/ce34d73f1675852d8f8b064e9f312e12a0fc9b10
RUN git checkout ce34d73f1675852d8f8b064e9f312e12a0fc9b10
RUN lein deps
RUN lein compile

ENTRYPOINT ["tini", "--"]

COPY --from=lsp-adapter /usr/local/bin/lsp-adapter /usr/local/bin/
EXPOSE 8080
# clojure-lsp does not support int ID fields https://github.com/snoe/clojure-lsp/issues/4
CMD ["lsp-adapter", "-proxyAddress=0.0.0.0:8080", "-jsonrpc2IDRewrite=string", "lein", "run"]

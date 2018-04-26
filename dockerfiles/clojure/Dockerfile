FROM golang:1.10-alpine
WORKDIR /go/src/github.com/sourcegraph/lsp-adapter
COPY . .
RUN CGO_ENABLED=0 GOBIN=/usr/local/bin go install github.com/sourcegraph/lsp-adapter

FROM clojure:lein-2.8.1-alpine
RUN apk add --no-cache ca-certificates git
RUN git clone https://github.com/snoe/clojure-lsp.git
WORKDIR clojure-lsp
RUN lein deps
RUN lein compile

COPY --from=0 /usr/local/bin/lsp-adapter /usr/local/bin/
EXPOSE 8080
# clojure-lsp does not support int ID fields https://github.com/snoe/clojure-lsp/issues/4
CMD ["lsp-adapter", "-proxyAddress=0.0.0.0:8080", "-jsonrpc2IDRewrite=string", "lein", "run"]

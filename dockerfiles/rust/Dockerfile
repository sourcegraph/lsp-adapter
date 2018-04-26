FROM golang:1.10-alpine
WORKDIR /go/src/github.com/sourcegraph/lsp-adapter
COPY . .
RUN CGO_ENABLED=0 GOBIN=/usr/local/bin go install github.com/sourcegraph/lsp-adapter

FROM rust:jessie
RUN rustup update
RUN rustup component add rls-preview rust-analysis rust-src

COPY --from=0 /usr/local/bin/lsp-adapter /usr/local/bin/
EXPOSE 8080
# RLS does not support string ID fields https://github.com/rust-lang-nursery/rls/issues/805
CMD ["lsp-adapter", "-proxyAddress=0.0.0.0:8080", "-jsonrpc2IDRewrite=number", "rls"]

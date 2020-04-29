FROM golang:1.14-alpine AS lsp-adapter
WORKDIR /go/src/github.com/sourcegraph/lsp-adapter
COPY . .
RUN CGO_ENABLED=0 GOBIN=/usr/local/bin go install github.com/sourcegraph/lsp-adapter

FROM rust:jessie
RUN rustup default nightly
RUN rustup update
RUN rustup component add rls-preview rust-analysis rust-src

ENV TINI_VERSION v0.18.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]

COPY --from=lsp-adapter /usr/local/bin/lsp-adapter /usr/local/bin/
EXPOSE 8080
# RLS does not support string ID fields https://github.com/rust-lang-nursery/rls/issues/805
CMD ["lsp-adapter", "-trace", "-proxyAddress=0.0.0.0:8080", "-jsonrpc2IDRewrite=number", "rls"]

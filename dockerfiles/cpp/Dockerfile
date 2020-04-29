FROM golang:1.14-alpine AS lsp-adapter
WORKDIR /go/src/github.com/sourcegraph/lsp-adapter
COPY . .
RUN CGO_ENABLED=0 GOBIN=/usr/local/bin go install github.com/sourcegraph/lsp-adapter

FROM debian:jessie

# Use tini as entrypoint to correctly handle signals and zombie processes.
ENV TINI_VERSION v0.18.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]

# See https://storage.googleapis.com/vim-clangd/${LLVM_COMMIT}/clangd-download-index.json
# Thank you https://github.com/Chilledheart/vim-clangd
ENV LLVM_COMMIT 796005d4a195f37222d7f63dc38b7b8ffe8cef6d
ADD https://storage.googleapis.com/vim-clangd/${LLVM_COMMIT}/clangd-debian-8.tar.gz .
RUN tar xf clangd-debian-8.tar.gz -C /usr/local \
  && rm clangd-debian-8.tar.gz \
  && chmod +x /usr/local/bin/clangd

COPY --from=lsp-adapter /usr/local/bin/lsp-adapter /usr/local/bin
EXPOSE 8080
CMD ["lsp-adapter", "--proxyAddress=0.0.0.0:8080", "-didOpenLanguage=cpp", "clangd"]

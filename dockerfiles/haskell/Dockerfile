FROM golang:1.14-alpine AS lsp-adapter
WORKDIR /go/src/github.com/sourcegraph/lsp-adapter
COPY . .
RUN CGO_ENABLED=0 GOBIN=/usr/local/bin go install github.com/sourcegraph/lsp-adapter

# sourcegraph/haskell-ide-engine is an image with a compiled `hie` executable.
# It has to be a separate, manually-built, image because the Travis environment
# only provides 3 GB of RAM, which is not enough to compile haskell-ide-engine.
# The Dockerfile is defined in ../haskell-ide-engine/Dockerfile.
FROM sourcegraph/haskell-ide-engine:562ac94

ENV TINI_VERSION v0.18.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]

COPY --from=lsp-adapter /usr/local/bin/lsp-adapter /usr/local/bin/
EXPOSE 8080
# -didOpenLanguage is necessary, otherwise haskell-ide-engine won't build the
# project.
CMD ["lsp-adapter", "-trace", "-proxyAddress=0.0.0.0:8080", "-didOpenLanguage=haskell", "/root/.local/bin/hie", "--lsp"]

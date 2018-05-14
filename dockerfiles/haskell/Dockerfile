FROM golang:1.10-alpine
WORKDIR /go/src/github.com/sourcegraph/lsp-adapter
COPY . .
RUN CGO_ENABLED=0 GOBIN=/usr/local/bin go install github.com/sourcegraph/lsp-adapter

FROM ubuntu:16.04
RUN apt-get update \
  # libtinfo-dev is required by haskell-ide-engine.
  && apt-get install -y git wget libtinfo-dev \
  && wget -qO- https://get.haskellstack.org/ | sh \
  && git clone https://github.com/haskell/haskell-ide-engine --recursive \
  && cd haskell-ide-engine \
  && git checkout 753deb91f9b9b39f14722315277d9fd587716bde \
  # haskell-ide-engine only supports the version of the Haskell compiler (GHC)
  # that it was compiled with. 8.2.2 is fairly recent (November 2017) and is
  # shown as an example in the README:
  # https://github.com/haskell/haskell-ide-engine/tree/753deb91f9b9b39f14722315277d9fd587716bde#for-ghc-822
  && stack --stack-yaml=stack-8.2.2.yaml install \
  && find . -name .git | xargs rm -rf

ENV TINI_VERSION v0.18.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]

COPY --from=0 /usr/local/bin/lsp-adapter /usr/local/bin/
EXPOSE 8080
# -didOpenLanguage is necessary, otherwise haskell-ide-engine won't build the
# project.
CMD ["lsp-adapter", "-trace", "-proxyAddress=0.0.0.0:8080", "-didOpenLanguage=haskell", "/root/.local/bin/hie", "--lsp"]

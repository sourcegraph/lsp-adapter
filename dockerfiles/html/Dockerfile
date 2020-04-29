FROM golang:1.14-alpine AS lsp-adapter
WORKDIR /go/src/github.com/sourcegraph/lsp-adapter
COPY . .
RUN CGO_ENABLED=0 GOBIN=/usr/local/bin go install github.com/sourcegraph/lsp-adapter

FROM node:9

# Use tini as entrypoint to correctly handle signals and zombie processes.
ENV TINI_VERSION v0.18.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]

USER node
# see https://github.com/nodejs/docker-node/blob/master/docs/BestPractices.md#global-npm-dependencies
ENV NPM_CONFIG_PREFIX=/home/node/.npm-global
ENV PATH=$PATH:/home/node/.npm-global/bin
# Ignore the warning about using WORKDIR instead of cd for convenience.
# hadolint ignore=DL3003
RUN cd "$HOME" \
  && git clone https://github.com/chrismwendt/vscode-html-languageserver-bin --recursive \
  && cd vscode-html-languageserver-bin \
  && git checkout 5e34582126c02387b490d7cbadacc45020bc70f4 \
  && npm install \
  && npm run build \
  && cd dist \
  && npm install --global

COPY --from=lsp-adapter /usr/local/bin/lsp-adapter /usr/local/bin
EXPOSE 8080
# Modify this command to connect to the language server
CMD ["lsp-adapter", "--trace", "--didOpenLanguage='html'", "--glob=*.html", "--proxyAddress=0.0.0.0:8080", "html-languageserver", "--stdio"]

FROM golang:1.14-alpine AS lsp-adapter
WORKDIR /go/src/github.com/sourcegraph/lsp-adapter
COPY . .
RUN CGO_ENABLED=0 GOBIN=/usr/local/bin go install github.com/sourcegraph/lsp-adapter

# 👀 Add steps here to build the language server itself 👀
# CMD ["echo", "🚨 This statement should be removed once you have added the logic to start up the language server! 🚨 Exiting..."]

FROM rocker/rstudio:3.5.1

RUN R -e 'chooseCRANmirror(graphics=FALSE, ind=65);source("https://install-github.me/REditorSupport/languageserver")'

ENV TINI_VERSION v0.18.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]

# Modify these commands to connect to the language server
COPY --from=lsp-adapter /usr/local/bin/lsp-adapter /usr/local/bin
EXPOSE 8080
CMD ["lsp-adapter", "--didOpenLanguage=r", "--trace", "--proxyAddress=0.0.0.0:8080", "R", "--quiet", "--slave", "-e", "languageserver::run()"]

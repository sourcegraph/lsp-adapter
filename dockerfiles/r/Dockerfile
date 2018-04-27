FROM golang:1.10-alpine
WORKDIR /go/src/github.com/sourcegraph/lsp-adapter
COPY . .
RUN CGO_ENABLED=0 GOBIN=/usr/local/bin go install github.com/sourcegraph/lsp-adapter

# 👀 Add steps here to build the language server itself 👀
# CMD ["echo", "🚨 This statement should be removed once you have added the logic to start up the language server! 🚨 Exiting..."]

FROM rocker/rstudio:3.4.4

RUN R -e 'chooseCRANmirror(graphics=FALSE, ind=65);source("https://install-github.me/REditorSupport/languageserver")'

# Modify these commands to connect to the language server
COPY --from=0 /usr/local/bin/lsp-adapter .
EXPOSE 8080
ENTRYPOINT ./lsp-adapter --didOpenLanguage='r' --trace --proxyAddress=0.0.0.0:8080 R --quiet --slave -e 'languageserver::run()'

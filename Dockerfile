FROM golang:1.10-alpine
WORKDIR /go/src/github.com/sourcegraph/lsp-adapter
COPY . .
RUN CGO_ENABLED=0 GOBIN=/usr/local/bin go install github.com/sourcegraph/lsp-adapter

# 👀 Add steps here to build the language server itself 👀
CMD ["echo", "🚨 This statement should be removed once you have added the logic to start up the language server! 🚨 Exiting..."]

# Modify these commands to connect to the language server
#COPY --from=0 /usr/local/bin/lsp-adapter .
#EXPOSE 8080
#CMD ["./lsp-adapter", "--proxyAddress=0.0.0.0:8080", ...]

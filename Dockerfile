FROM golang:1.14-alpine
WORKDIR /go/src/github.com/sourcegraph/lsp-adapter
COPY . .
RUN CGO_ENABLED=0 GOBIN=/usr/local/bin go install github.com/sourcegraph/lsp-adapter

# 👀 Add steps here to build the language server itself 👀
CMD ["echo", "🚨 This statement should be removed once you have added the logic to start up the language server! 🚨 Exiting..."]

# # Use tini as entrypoint to correctly handle signals and zombie processes.
# ENV TINI_VERSION v0.18.0
# ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
# RUN chmod +x /tini
# ENTRYPOINT ["/tini", "--"]

# COPY --from=0 /usr/local/bin/lsp-adapter /usr/local/bin
# EXPOSE 8080
# # Modify this command to connect to the language server
# CMD ["lsp-adapter", "--proxyAddress=0.0.0.0:8080", ...]

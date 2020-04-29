FROM golang:1.14-alpine AS lsp-adapter
WORKDIR /go/src/github.com/sourcegraph/lsp-adapter
COPY . .
RUN CGO_ENABLED=0 GOBIN=/usr/local/bin go install github.com/sourcegraph/lsp-adapter

# 👀 Add steps here to build the language server itself 👀
# CMD ["echo", "🚨 This statement should be removed once you have added the logic to start up the language server! 🚨 Exiting..."]

FROM elixir:1.6

# Use tini as entrypoint to correctly handle signals and zombie processes.
ENV TINI_VERSION v0.18.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]

ENV EXLIXIR_LS_VERSION v0.2.18
ADD https://github.com/JakeBecker/elixir-ls/releases/download/${EXLIXIR_LS_VERSION}/elixir-ls.zip .
RUN apt-get update && apt-get install --no-install-recommends -y unzip && rm -rf /var/lib/apt/lists/*
RUN unzip elixir-ls.zip && rm elixir-ls.zip
RUN chmod +x language_server.sh

COPY --from=lsp-adapter /usr/local/bin/lsp-adapter /usr/local/bin
EXPOSE 8080
# Modify this command to connect to the language server
CMD ["lsp-adapter", "--trace", "--didOpenLanguage='elixir'", "--proxyAddress=0.0.0.0:8080", "./language_server.sh"]

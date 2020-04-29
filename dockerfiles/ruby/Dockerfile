FROM golang:1.14-alpine AS lsp-adapter
WORKDIR /go/src/github.com/sourcegraph/lsp-adapter
COPY . .
RUN CGO_ENABLED=0 GOBIN=/usr/local/bin go install github.com/sourcegraph/lsp-adapter
RUN CGO_ENABLED=0 go build -o /usr/local/bin/solargraph-wrapper dockerfiles/ruby/solargraph-wrapper.go

FROM ruby:2.5

RUN gem install solargraph

ENV TINI_VERSION v0.18.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]

COPY --from=lsp-adapter /usr/local/bin/lsp-adapter /usr/local/bin/
COPY --from=lsp-adapter /usr/local/bin/solargraph-wrapper /usr/local/bin/
EXPOSE 8080
# Solargraph issues:
# * Expects string to be a number
# * Expects didOpen to be sent
# * Only supports TCP connections (solargraph-wrapper)
# * Requires CWD to be rootURI (solargraph-wrapper)
CMD ["lsp-adapter", "-proxyAddress=0.0.0.0:8080", "-jsonrpc2IDRewrite=string", "-didOpenLanguage=ruby", "solargraph-wrapper"]

FROM golang:1.10-alpine
WORKDIR /go/src/github.com/sourcegraph/lsp-adapter
COPY . .
RUN CGO_ENABLED=0 GOBIN=/usr/local/bin go install github.com/sourcegraph/lsp-adapter

FROM ruby:2.5
#RUN apt-get update && apt-get install -y netcat && rm -rf /var/lib/apt/lists/*

RUN gem install solargraph

COPY --from=0 /usr/local/bin/lsp-adapter /usr/local/bin/
EXPOSE 8080
CMD ["lsp-adapter", "-proxyAddress=0.0.0.0:8080", "-jsonrpc2IDRewrite=number"]

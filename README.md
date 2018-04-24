# lsp-adapter [![Travis-CI](https://travis-ci.org/sourcegraph/lsp-adapter.svg)](https://travis-ci.org/sourcegraph/lsp-adapter) [![AppVeyor](https://ci.appveyor.com/api/projects/status/vfdftcqh0ekb881u/branch/master?svg=true)](https://ci.appveyor.com/project/sourcegraph/lsp-adapter/branch/master) [![GoDoc](https://godoc.org/github.com/sourcegraph/lsp-adapter?status.svg)](http://godoc.org/github.com/sourcegraph/lsp-adapter) [![Report card](https://goreportcard.com/badge/github.com/sourcegraph/lsp-adapter)](https://goreportcard.com/report/github.com/sourcegraph/lsp-adapter)

Command `lsp-adapter` provides a proxy which adapts Sourcegraph LSP requests to
vanilla LSP requests.

## Background

[Code Intelligence on Sourcegraph](https://about.sourcegraph.com/docs/code-intelligence/) is powered by the [Language Server Protocol](https://microsoft.github.io/language-server-protocol/). 

Previously, language servers that were used on sourcegraph.com were additionally required to support our
custom LSP [files extensions](https://github.com/sourcegraph/language-server-protocol/blob/master/extension-files.md). These extensions allowed language servers to operate without sharing a physical file system with the client. While it's preferable for language servers to implement these extensions for performance reasons, implementing this functionality is a large undertaking.  

`lsp-adapter` eliminates the need for this requirement, which allows off-the-shelf language servers to be able to provide basic functionality (hovers, local definitions) to Sourcegraph.

## How to Install 

To build and install `lsp-adapter`, run:

```shell
go get -u -v github.com/sourcegraph/lsp-adapter
```

Running `lsp-adapter --help` shows you some of its options: 

```shell
> lsp-adapter --help                                                                                          
Usage: lsp-adapter [OPTIONS] LSP_COMMAND_ARGS...

Options:
  -cacheDirectory string
    	cache directory location (default "/var/folders/qq/1q_cmsmx6qv7bs_m6g_2pt1r0000gn/T/proxy-cache")
  -didOpenLanguage string
    	(HACK) If non-empty, send 'textDocument/didOpen' notifications with the specified language field (e.x. 'python') to the language server for every file.
  -proxyAddress string
    	proxy server listen address (tcp) (default "127.0.0.1:8080")
  -trace
    	trace logs to stderr
```
## How to Use `lsp-adapter`

`lsp-adapter` proxies requests between your Sourcegraph instance and the language server, and modifies them in such a way that allows for the two to communicate correctly. In order to do this, we need to know

- How to connect to connect `lsp-adapter` to the language server
- How to connect to your Sourcegraph instance to `lsp-adapter`


### Connect `lsp-adapter` to the Language Server 

`lsp-adapter` can talk to language servers over standard I/O.

`lsp-adapter` interprets any positional arguments after the flags as the necessary command (+ arguments) to start the language server binary. It uses this command to communicate to the language server inside of a subprocess.

For example, if I am trying to use Rust’s language server, the command to start it up is just rls. `lsp-adapter` can be told to start the start the same server via:

```shell
lsp-adapter rls
```

Any stderr output from the binary will also appear in `lsp-adapter`'s logs.


### Connect Sourcegraph to `lsp-adapter`

1. Use the `-proxyAddress` flag to tell `lsp-adapter` what address to listen for connections from Sourcegraph on. For example, I can tell `lsp-adapter` to listen on my local `8080` port with `-proxyAddress=127.0.0.1:8080`.

2. We then need to add a new entry to the `"langservers"` field in the site configuration in order to point Sourcegraph at `lsp-adapter` (similar to the steps in [this document](https://about.sourcegraph.com/docs/code-intelligence/install-manual)). For example, if `lsp-adapter` is connected to the Rust language server, and the `lsp-adapter` itself is listening on `127.0.0.1:8080`:

```json
{
    "language": "rust",
    "address": "tcp://127.0.0.1:8080"
}
```
would be the new entry that needs to be added to the `"langservers"` field.

## Example Commands 

Connect via standard I/O to a language server whose command can be run with `rls`, and listen for connections from Sourcegraph from any address on port `1234`.

```shell
> lsp-adapter -proxyAddresss=0.0.0.0:1234 rls

2018/04/04 15:27:04 proxy.go:71: CloneProxy: accepting connections at [::]:8080
```

Connect via standard I/O to a language server whose command can be run with `rls`, change the location of the cache directory (used for cloning the repo locally) to `/tmpDir`, and enable tracing for every request to/from the language server.

```shell
> lsp-adapter -cacheDir='/tmpDir' -trace rls
```

## Docker 

There is a [skeleton Dockerfile](https://github.com/sourcegraph/lsp-adapter/blob/master/Dockerfile) that shows how to package `lsp-adapter` along with your desired language server inside of a docker container. 

A fully working example for Rust:

```Dockerfile
FROM golang:1.10-alpine
WORKDIR /lsp-adapter
RUN apk add --no-cache ca-certificates git
COPY . .
RUN CGO_ENABLED=0 go get -d -v ./...
RUN CGO_ENABLED=0 go build -o lsp-adapter *.go

# 👀 Add steps here to build the language server itself 👀
# CMD ["echo", "🚨 This statement should be removed once you have added the logic to start up the language server! 🚨 Exiting..."]

FROM rust:jessie
RUN rustup update
RUN rustup component add rls-preview rust-analysis rust-src

# Modify these commands to connect to the language server
COPY --from=0 /lsp-adapter/lsp-adapter .
EXPOSE 8080
CMD ["./lsp-adapter", "--proxyAddress=0.0.0.0:8080", "rls"]
```

## Did Open Hack 

Some language servers incorrectly follow the LSP spec, and refuse to work unless the `textDocument/didOpen` notification has been sent. See [this commit](https://github.com/sourcegraph/lsp-adapter/commit/1228a1fbaf102aa44575cec6802a5a211d117ee1) for more context. If the language server that you’re trying to use has this issue, try setting the `didOpenLanguage` flag (example: if a python language server had this issue - use `./lsp-adapter -didOpenLanguage=’python’...`) to work around it. 
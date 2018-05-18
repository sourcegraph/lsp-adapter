# lsp-adapter [![Travis-CI](https://travis-ci.org/sourcegraph/lsp-adapter.svg)](https://travis-ci.org/sourcegraph/lsp-adapter) [![AppVeyor](https://ci.appveyor.com/api/projects/status/vfdftcqh0ekb881u/branch/master?svg=true)](https://ci.appveyor.com/project/sourcegraph/lsp-adapter/branch/master) [![GoDoc](https://godoc.org/github.com/sourcegraph/lsp-adapter?status.svg)](http://godoc.org/github.com/sourcegraph/lsp-adapter) [![Go Report Card](https://goreportcard.com/badge/github.com/sourcegraph/lsp-adapter)](https://goreportcard.com/report/github.com/sourcegraph/lsp-adapter)

Command `lsp-adapter` provides a proxy which adapts Sourcegraph LSP requests to
vanilla LSP requests.

## Background

[Code Intelligence on Sourcegraph](https://about.sourcegraph.com/docs/code-intelligence/) is powered by the [Language Server Protocol](https://microsoft.github.io/language-server-protocol/).

Previously, language servers that were used on sourcegraph.com were additionally required to support our
custom LSP [files extensions](https://github.com/sourcegraph/language-server-protocol/blob/master/extension-files.md). These extensions allowed language servers to operate without sharing a physical file system with the client. While it's preferable for language servers to implement these extensions for performance reasons, implementing this functionality is a large undertaking.

`lsp-adapter` eliminates the need for this requirement, which allows off-the-shelf language servers to be able to provide basic functionality (hovers, local definitions) to Sourcegraph.

## How to Install

You can download the latest binary on the [releases page](https://github.com/sourcegraph/lsp-adapter/releases).

Alternatively, install it with `go get`:

```shell
go get -u -v github.com/sourcegraph/lsp-adapter
```

Running `lsp-adapter --help` shows you some of its options:

```shell
> lsp-adapter --help
Usage: lsp-adapter [OPTIONS] LSP_COMMAND_ARGS...

Options:
  -cacheDirectory string
      cache directory location (default "/tmp/proxy-cache")
  -didOpenLanguage string
      (HACK) If non-empty, send 'textDocument/didOpen' notifications with the specified language field (e.x. 'python') to the language server for every file.
  -glob string
      A colon (:) separated list of file globs to sync locally. By default we place all files into the workspace, but some language servers may only look at a subset of files. Specifying this allows us to avoid syncing all files. Note: This is done base name only.
  -jsonrpc2IDRewrite string
      (HACK) Rewrite jsonrpc2 ID. none (default) is no rewriting. string will use a string ID. number will use a number ID. Useful for language servers with non-spec complaint JSONRPC2 implementations. (default "none")
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

There is a [skeleton Dockerfile](./Dockerfile) that shows how to package `lsp-adapter` along with your desired language server inside of a docker container. There are fully working examples in [dockerfiles](./dockerfiles). For example [dockerfiles/rust/Dockerfile](./dockerfiles/rust/Dockerfile), which can be built with:

```shell
> docker build -f dockerfiles/rust/Dockerfile .
```

## Glob

Most language servers will only ever look at files that match a set of known patterns. On initialize lsp-adapter copies a full work-tree to disk for a repository, but by specify `-glob` we can avoid copying over files that will not be looked at. For example, if a python language server only looks at `py` and `pyc` files you can specify `-glob=*.py:*.pyc`. The matching is done on the basename of the path using [path.Match](https://godoc.org/path#Match).

## Did Open Hack

Some language servers do not follow the LSP spec correctly and refuse to work unless the `textDocument/didOpen` notification has been sent. See [this commit](https://github.com/sourcegraph/lsp-adapter/commit/1228a1fbaf102aa44575cec6802a5a211d117ee1) for more context. If the language server that you’re trying to use has this issue, try setting the `didOpenLanguage` flag (example: if a python language server had this issue - use `./lsp-adapter -didOpenLanguage=python ...`) to work around it.


## JSONRPC2 ID Rewrite Hack

Some language servers do not follow the JSONRPC2 spec correctly and fail if the Request ID is not a number of string. If the language server that you’re trying to use has this issue, try setting the `jsonrpc2IDRewrite` flag (example: if a rust language server had this issue - use `./lsp-adapter -jsonrpc2IDRewrite=number ...`) to work around it.

# Experimental Haskell language server

![demonstration of hovers from the Haskell language server inside chrismwendt/minijava](https://cl.ly/1h3E2P2s2g2y/haskell.gif)

*This GIF was created by browsing http://localhost:3080/github.com/chrismwendt/MiniJava@6cb615856cfcad0253c9588a40a5b8678df05349/-/blob/src/RegisterAllocator.hs with this Haskell language server enabled.*

## Introduction

This Dockerfile adds experimental Haskell language support for Sourcegraph.

Thanks to the [haskell/haskell-ide-engine](https://github.com/haskell/haskell-ide-engine) project for providing the language server that's wrapped by `lsp-adapter` in this image.

Check out the [Sourcegraph docs](http://about.sourcegraph.com/docs/code-intelligence/experimental-language-servers) for information on enabling this language server for your Sourcegraph installation.

## How to update the language server

The [sourcegraph/codeintel-haskell](https://hub.docker.com/r/sourcegraph/codeintel-haskell) image follows a different build process from the other images because Travis is unable to compile [haskell-ide-engine](https://github.com/haskell/haskell-ide-engine) due to limited RAM.

The way around this is to build an intermediate image [sourcegraph/haskell-ide-engine](https://hub.docker.com/r/sourcegraph/haskell-ide-engine) elsewhere (e.g. on your own machine), then have Travis use that as a base on which to build [sourcegraph/codeintel-haskell](https://hub.docker.com/r/sourcegraph/codeintel-haskell).

In order to update the version of [haskell-ide-engine](https://github.com/haskell/haskell-ide-engine) in [sourcegraph/codeintel-haskell](https://hub.docker.com/r/sourcegraph/codeintel-haskell):

1. Follow [../haskell-ide-engine/README.md](../haskell-ide-engine/README.md) to update [../haskell-ide-engine/Dockerfile](../haskell-ide-engine/Dockerfile) then build and push [sourcegraph/haskell-ide-engine](https://hub.docker.com/r/sourcegraph/haskell-ide-engine).
1. In [./Dockerfile](./Dockerfile), replace the tag in the line `FROM sourcegraph/haskell-ide-engine:562ac94` with the tag you just created.
1. Commit **both** of the Dockerfiles and submit a PR.
1. After the PR is merged, Travis will automatically deploy [sourcegraph/codeintel-haskell](https://hub.docker.com/r/sourcegraph/codeintel-haskell).

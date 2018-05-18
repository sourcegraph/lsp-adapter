# Experimental Haskell language server

![demonstration of hovers from the Haskell language server inside chrismwendt/minijava](https://cl.ly/1h3E2P2s2g2y/haskell.gif)

*This GIF was created by browsing http://localhost:3080/github.com/chrismwendt/MiniJava@6cb615856cfcad0253c9588a40a5b8678df05349/-/blob/src/RegisterAllocator.hs with this Haskell language server enabled.*

## Introduction

This Dockerfile adds experimental Haskell language support for Sourcegraph.

Thanks to the [haskell/haskell-ide-engine](https://github.com/haskell/haskell-ide-engine) project for providing the language server that's wrapped by `lsp-adapter` in this image.

Check out the [Sourcegraph docs](http://about.sourcegraph.com/docs/code-intelligence/experimental-language-servers) for information on enabling this language server for your Sourcegraph installation.
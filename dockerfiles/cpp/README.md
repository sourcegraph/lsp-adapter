# Experimental C++ language server 

![demonstration of hovers from the C++ language server inside antirez/redis](https://cl.ly/0f401u080U3S/Screen%20Recording%202018-04-30%20at%2012.21%20pm.gif)

*This GIF was created by browsing https://sourcegraph.com/github.com/antirez/redis/-/blob/src/expire.c with this C++ language server enabled.*

## Introduction

This Dockerfile adds experimental C++ language support for Sourcegraph. 

Thanks to [LLVM](https://clang.llvm.org/extra/clangd.html) (and [Chilledheart/vim-clangd](https://github.com/Chilledheart/vim-clangd) for providing prebuilt binaries) for providing the language server that's wrapped by `lsp-adapter` in this image.

Check out the [Sourcegraph docs](http://about.sourcegraph.com/docs/code-intelligence/experimental-language-servers) for information on enabling this language server for your Sourcegraph installation.

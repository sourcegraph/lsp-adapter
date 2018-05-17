# Experimental HTML language server 

![demonstration of hovers from the HTML language server inside tabler/tabler](https://cl.ly/1W2a1u2n3H0N/Screen%20Recording%202018-05-07%20at%2005.23%20PM.gif)

*This GIF was created by browsing https://sourcegraph.com/github.com/tabler/tabler/-/blob/dist/login.html with this HTML language server enabled.*

## Introduction

This Dockerfile adds experimental HTML language support for Sourcegraph. 

Thanks to the [Microsoft/vscode](https://github.com/Microsoft/vscode/) project (and [vscode-langservers/vscode-html-languageserver-bin](https://github.com/vscode-langservers/vscode-html-languageserver-bin) for providing prebuilt binaries) for providing the language server that's wrapped by `lsp-adapter` in this image.

Check out the [Sourcegraph docs](http://about.sourcegraph.com/docs/code-intelligence/experimental-language-servers) for information on enabling this language server for your Sourcegraph installation.

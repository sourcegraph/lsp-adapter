# Experimental CSS language server 

![demonstration of hovers from the CSS language server inside tabler/tabler](https://cl.ly/3I2S3x2Z3K44/Screen%20Recording%202018-05-07%20at%2004.36%20PM.gif)

*This GIF was created by browsing https://sourcegraph.com/github.com/tabler/tabler/-/blob/dist/assets/plugins/prismjs/plugin.css with this CSS language server enabled.*

## Introduction

This Dockerfile adds experimental CSS language support for Sourcegraph. 

Thanks to the [Microsoft/vscode](https://github.com/Microsoft/vscode/) project (and [vscode-langservers/vscode-css-languageserver-bin](https://github.com/vscode-langservers/vscode-css-languageserver-bin) for providing prebuilt binaries) for providing the language server that's wrapped by `lsp-adapter` in this image.

Check out the [Sourcegraph docs](http://about.sourcegraph.com/docs/code-intelligence/experimental-language-servers) for information on enabling this language server for your Sourcegraph installation.

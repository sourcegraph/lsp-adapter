# Experimental Bash language server 

![demonstration of hovers from the Bash language server inside creationix/nvm](https://cl.ly/271p292i342p/Screen%20Recording%202018-05-07%20at%2009.19%20AM.gif)

*This GIF was created by browsing https://sourcegraph.com/github.com/creationix/nvm/-/blob/nvm.sh with this Bash language server enabled.*

## Introduction

This Dockerfile adds experimental Bash language support for Sourcegraph. 

Thanks to the [mads-hartmann/bash-language-server](https://github.com/mads-hartmann/bash-language-server) project for providing the language server that's wrapped by `lsp-adapter` in this image.

Check out the [Sourcegraph docs](http://about.sourcegraph.com/docs/code-intelligence/experimental-language-servers) for information on enabling this language server for your Sourcegraph installation.

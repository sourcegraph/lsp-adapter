# Experimental Dockerfile language server 

![demonstration of hovers from the Dockerfile language server inside moby/datakit](https://cl.ly/2f3P2t2j2W2i/Screen%20Recording%202018-05-07%20at%2004.55%20PM.gif)

<!-- TODO(@ggilmore @keegancsmith @felixfbecker): Revisit creating this GIF once improved tooltip and syntax highlighting code lands-->

*This GIF was created by browsing https://sourcegraph.com/github.com/moby/datakit/-/blob/ci/self-ci/Dockerfile with this Dockerfile language server enabled.*

## Introduction

This Dockerfile adds experimental Dockerfile language support for Sourcegraph. 

Thanks to the [rcjsuen/dockerfile-language-server-nodejs](https://github.com/rcjsuen/dockerfile-language-server-nodejs) project for providing the language server that's wrapped by `lsp-adapter` in this image.

Check out the [Sourcegraph docs](http://about.sourcegraph.com/docs/code-intelligence/experimental-languages) for information on enabling this language server for your Sourcegraph installation.

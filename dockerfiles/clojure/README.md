# Experimental Clojure language server 

![demonstration of hovers from the Clojure language server inside ryangreenhall/clojure-hello-world$](https://cl.ly/001Z380o0K0B/Screen%20Recording%202018-05-07%20at%2004.09%20PM.gif)

*This GIF was created by browsing https://sourcegraph.com/github.com/ryangreenhall/clojure-hello-world/-/blob/src/clojure_hello_world/core.clj with this Clojure language server enabled.*

## Introduction

This Dockerfile adds experimental Clojure language support for Sourcegraph. 

Thanks to the [snoe/clojure-lsp](https://github.com/snoe/clojure-lsp) project for providing the language server that's wrapped by `lsp-adapter` in this image.

Check out the [Sourcegraph docs](http://about.sourcegraph.com/docs/code-intelligence/preview-languages) for information on enabling this language server for your Sourcegraph installation.

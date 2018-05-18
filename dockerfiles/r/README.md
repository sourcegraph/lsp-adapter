# Experimental R language server

![demonstration of hovers from the R language server inside blmoore/mandelbrot](https://cl.ly/2G2k1G0V3P3C/r.gif)

*This GIF was created by browsing https://sourcegraph.com/github.com/blmoore/mandelbrot/-/blob/R/util.R with this R language server enabled.*

## Introduction

This Dockerfile adds experimental R language support for Sourcegraph.

Thanks to the [REditorSupport/languageserver](https://github.com/REditorSupport/languageserver) project for providing the language server that's wrapped by `lsp-adapter` in this image.

Check out the [Sourcegraph docs](http://about.sourcegraph.com/docs/code-intelligence/experimental-language-servers) for information on enabling this language server for your Sourcegraph installation.

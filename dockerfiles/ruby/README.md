# Experimental Ruby language server 

![demonstration of hovers from the Ruby language server inside tabler/tabler](https://cl.ly/3l0Z2g1l3D2r/Screen%20Recording%202018-05-07%20at%2005.30%20PM.gif)

*This GIF was created by browsing https://sourcegraph.com/github.com/tabler/tabler/-/blob/src/_plugins/jekyll-tabler.rb with this Ruby language server enabled.*

## Introduction

This Dockerfile adds experimental Ruby language support for Sourcegraph. 

Thanks to the [castwide/solargraph](https://github.com/castwide/solargraph) project for providing the language server that's wrapped by `lsp-adapter` in this image.

Check out the [Sourcegraph docs](http://about.sourcegraph.com/docs/code-intelligence/experimental-language-servers) for information on enabling this language server for your Sourcegraph installation.

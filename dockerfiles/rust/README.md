# Experimental Rust language server 

![demonstration of hovers from the Rust language server inside databricks/click](https://cl.ly/383f3V0P1r1u/Screen%20Recording%202018-05-07%20at%2005.58%20PM.gif)

*This GIF was created by browsing https://sourcegraph.com/github.com/databricks/click/-/blob/src/describe.rs with this Rust language server enabled.*

## Introduction

This Dockerfile adds experimental Rust language support for Sourcegraph. 

Thanks to the [rust-lang-nursery/rls](https://github.com/rust-lang-nursery/rls) project for providing the language server that's wrapped by `lsp-adapter` in this image.

Check out the [Sourcegraph docs](http://about.sourcegraph.com/docs/code-intelligence/preview-languages) for information on enabling this language server for your Sourcegraph installation.

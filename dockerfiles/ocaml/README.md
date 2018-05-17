# Experimental OCaml language server 

![demonstration of hovers from the OCaml language server inside moby/datakit](https://cl.ly/3l0Z2g1l3D2r/Screen%20Recording%202018-05-07%20at%2005.30%20PM.gif)

*This GIF was created by browsing https://sourcegraph.com/github.com/moby/datakit/-/blob/ci/src/cI_docker.mli with this OCaml language server enabled.*

## Introduction

This Dockerfile adds experimental OCaml language support for Sourcegraph. 

Thanks to the [freebroccolo/ocaml-language-server](https://github.com/freebroccolo/ocaml-language-server) project for providing the language server that's wrapped by `lsp-adapter` in this image.

Check out the [Sourcegraph docs](http://about.sourcegraph.com/docs/code-intelligence/experimental-languages) for information on enabling this language server for your Sourcegraph installation.

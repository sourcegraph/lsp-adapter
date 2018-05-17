# Experimental Elixir language server 

![demonstration of hovers from the Elixir language server inside philnash/elixir-examples](https://cl.ly/3h3V0e3k3a21/Screen%20Recording%202018-05-07%20at%2005.11%20PM.gif)

<!-- TODO(@ggilmore @keegancsmith @felixfbecker): Revisit creating this GIF once improved tooltip and syntax highlighting code lands-->

*This GIF was created by browsing https://sourcegraph.com/github.com/philnash/elixir-examples/-/blob/hello-world/hello-world.exs with this Elixir language server enabled.*

## Introduction

This Dockerfile adds experimental Elixir language support for Sourcegraph. 

Thanks to the [JakeBecker/elixir-ls](https://github.com/JakeBecker/elixir-ls) project for providing the language server that's wrapped by `lsp-adapter` in this image.

Check out the [Sourcegraph docs](http://about.sourcegraph.com/docs/code-intelligence/experimental-language-servers) for information on enabling this language server for your Sourcegraph installation.

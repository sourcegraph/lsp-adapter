# Experimental Lua language server

![demonstration of hovers from the Lua language server inside daurnimator/luatz](https://cl.ly/2d1e08103r2A/lua.gif)

*This GIF was created by browsing https://sourcegraph.com/github.com/daurnimator/luatz/-/blob/luatz/gettime.lua with this Lua language server enabled.*

## Introduction

This Dockerfile adds experimental Lua language support for Sourcegraph.

Thanks to the [Alloyed/lua-lsp](https://github.com/Alloyed/lua-lsp) project providing the language server that's wrapped by `lsp-adapter` in this image.

Check out the [Sourcegraph docs](http://about.sourcegraph.com/docs/code-intelligence/experimental-language-servers) for information on enabling this language server for your Sourcegraph installation.

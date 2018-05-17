# Experimental C# language server

![](https://cl.ly/2R1f0D2e1I1w/csharp.gif)

*This GIF was created by browsing https://sourcegraph.com/github.com/adamreeve/semver.net@4890b46afd8c67a6eeb6a97a77a4adc24cf4d33a/-/blob/src/SemVer/Desugarer.cs with this C# language server enabled.*

## Introduction

This Dockerfile adds experimental C# language support for Sourcegraph.

Thanks to the [OmniSharp/omnisharp-node-client](https://github.com/OmniSharp/omnisharp-node-client) project for providing the language server that's wrapped by `lsp-adapter` in this image.

Check out the [Sourcegraph docs](https://about.sourcegraph.com/docs/code-intelligence/experimental-language-servers) for information on enabling this language server for your Sourcegraph installation.

# Experimental C# language server

![](https://cl.ly/2R1f0D2e1I1w/csharp.gif)

*This GIF was created by browsing https://sourcegraph.com/github.com/adamreeve/semver.net@4890b46afd8c67a6eeb6a97a77a4adc24cf4d33a/-/blob/src/SemVer/Desugarer.cs with this C# language server enabled.*

## Introduction

This Dockerfile adds experimental C# language support for Sourcegraph.

Thanks to the [OmniSharp/omnisharp-node-client](https://github.com/OmniSharp/omnisharp-node-client) project for providing the language server that's wrapped by `lsp-adapter` in this image.

## How to add this language server to Sourcegraph

*These steps are adapted from our documentation for [manually adding code intelligence to Sourcegraph server](https://about.sourcegraph.com/docs/code-intelligence/install-manual/).*

1. Run the `sourcegarph/server` Docker image:

```shell
docker run --publish 7080:7080 --rm --network=lsp --name=sourcegraph --volume ~/.sourcegraph/config:/etc/sourcegraph --volume ~/.sourcegraph/data:/var/opt/sourcegraph sourcegraph/server:2.7.6
```

2. Run the experimental C# language server:

  ```shell
  docker pull sourcegraph/codeintel-csharp

  docker run --rm --network=lsp --name=csharp sourcegraph/codeintel-csharp
  ```

3. Add the following entry to the `langservers` field in the [site configuration](https://about.sourcegraph.com/docs/config):

  ```json
  {
      // ...
      "langservers": [
          // ...
          {
              "language": "cs",
              "address": "tcp://csharp:8080",
              "metadata": {
                "experimental": true,
                "homepageURL": "https://github.com/OmniSharp/omnisharp-node-client",
                "issuesURL": "https://github.com/OmniSharp/omnisharp-node-client/issues",
                "docsURL": "https://github.com/OmniSharp/omnisharp-node-client/blob/master/README.md"
              }
          },
          // ...
      ]
      // ...
  }
  ```

**Note:** This language server lazily downloads OmniSharp once you open a C# file in Sourcegraph. You should see a message in the console once it's ready.
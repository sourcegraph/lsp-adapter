# Experimental Haskell language server

![demonstration of hovers from the Haskell language server inside chrismwendt/minijava](https://cl.ly/1h3E2P2s2g2y/haskell.gif)

*This GIF was created by browsing http://localhost:3080/github.com/chrismwendt/MiniJava@6cb615856cfcad0253c9588a40a5b8678df05349/-/blob/src/RegisterAllocator.hs with this Haskell language server enabled.*

## Introduction

This Dockerfile adds experimental Haskell language support for Sourcegraph.

Thanks to the [haskell/haskell-ide-engine](https://github.com/haskell/haskell-ide-engine) project for providing the language server that's wrapped by `lsp-adapter` in this image.

## How to add this language server to Sourcegraph

*These steps are adapted from our documentation for [manually adding code intelligence to Sourcegraph server](https://about.sourcegraph.com/docs/code-intelligence/install-manual/).*

1. Run the `sourcegarph/server` Docker image:

```shell
docker run --publish 7080:7080 --rm --network=lsp --name=sourcegraph --volume ~/.sourcegraph/config:/etc/sourcegraph --volume ~/.sourcegraph/data:/var/opt/sourcegraph sourcegraph/server:2.7.6
```

2. Run the experimental Dockerfile language server:

  ```shell
  docker pull sourcegraph/codeintel-haskell

  docker run --rm --network=lsp --name=haskell sourcegraph/codeintel-haskell
  ```

3. Add the following entry to the `langservers` field in the [site configuration](https://about.sourcegraph.com/docs/config):

  ```js
  {
      // ...
      "langservers": [
          // ...
          {
              "language": "haskell",
              "address": "tcp://haskell:8080",
              "metadata": {
                "experimental": true,
                "homepageURL": "https://github.com/haskell/haskell-ide-engine",
                "issuesURL": "https://github.com/haskell/haskell-ide-engine/issues",
                "docsURL": "https://github.com/haskell/haskell-ide-engine/blob/master/README.md"
              }
          },
          // ...
      ]
      // ...
  }
  ```

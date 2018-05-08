# Experimental Elixir language server 

![demonstration of hovers from the Elixir language server inside philnash/elixir-examples](https://cl.ly/3h3V0e3k3a21/Screen%20Recording%202018-05-07%20at%2005.11%20PM.gif)

<!-- TODO(@ggilmore @keegancsmith @felixfbecker): Revisit creating this GIF once improved tooltip and syntax highlighting code lands-->

*This GIF was created by browsing https://sourcegraph.com/github.com/philnash/elixir-examples/-/blob/hello-world/hello-world.exs with this Elixir language server enabled.*

## Introduction

This Dockerfile adds experimental Elixir language support for Sourcegraph. 

Thanks to the [JakeBecker/elixir-ls](https://github.com/JakeBecker/elixir-ls) project for providing the language server that's wrapped by `lsp-adapter` in this image.

## How to add this language server to Sourcegraph

*These steps are adapted from our documentation for [manually adding code intelligence to Sourcegraph server](https://about.sourcegraph.com/docs/code-intelligence/install-manual/).*

1. Run the `sourcegarph/server` Docker image: 

```shell
docker run --publish 7080:7080 --rm --network=lsp --name=sourcegraph --volume ~/.sourcegraph/config:/etc/sourcegraph --volume ~/.sourcegraph/data:/var/opt/sourcegraph sourcegraph/server:2.7.6
```

2. Run the experimental Elixir language server:

  ```shell
  docker pull sourcegraph/codeintel-elixir

  docker run --rm --network=lsp --name=elixir sourcegraph/codeintel-elixir
  ```

3. Add the following entry to the `langservers` field in the [site configuration](https://about.sourcegraph.com/docs/config):

  ```js
  {
      // ...
      "langservers": [
          // ...
          {
              "language": "elixir",
              "address": "tcp://elixir:8080",
              "metadata": {
                "experimental": true,
                "homepageURL": "https://github.com/JakeBecker/elixir-ls", 
                "issuesURL": "https://github.com/JakeBecker/elixir-ls/issues", 
                "docsURL": "https://github.com/JakeBecker/elixir-ls/blob/master/README.md"
              }
          },
          // ...
      ]
      // ...
  }
  ```

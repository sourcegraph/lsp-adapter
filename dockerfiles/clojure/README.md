# Experimental Clojure language server 

![demonstration of hovers from the Clojure language server inside ryangreenhall/clojure-hello-world$](https://cl.ly/001Z380o0K0B/Screen%20Recording%202018-05-07%20at%2004.09%20PM.gif)

*This GIF was created by browsing https://sourcegraph.com/github.com/ryangreenhall/clojure-hello-world/-/blob/src/clojure_hello_world/core.clj with this Clojure language server enabled.*

## Introduction

This Dockerfile adds experimental Clojure language support for Sourcegraph. 

Thanks to the [snoe/clojure-lsp](https://github.com/snoe/clojure-lsp) project for providing the language server that's wrapped by `lsp-adapter` in this image.

## How to add this language server to Sourcegraph

*These steps are adapted from our documentation for [manually adding code intelligence to Sourcegraph server](https://about.sourcegraph.com/docs/code-intelligence/install-manual/).*

1. Run the `sourcegarph/server` Docker image: 

```shell
docker run --publish 7080:7080 --rm --network=lsp --name=sourcegraph --volume ~/.sourcegraph/config:/etc/sourcegraph --volume ~/.sourcegraph/data:/var/opt/sourcegraph sourcegraph/server:2.7.6
```

2. Run the experimental Clojure language server:

  ```shell
  docker pull sourcegraph/codeintel-clojure

  docker run --rm --network=lsp --name=clojure sourcegraph/codeintel-clojure
  ```

3. Add the following entry to the `langservers` field in the [site configuration](https://about.sourcegraph.com/docs/config):

  ```js
  {
      // ...
      "langservers": [
          // ...
          {
              "language": "clojure",
              "address": "tcp://clojure:8080",
              "metadata": {
                "experimental": true,
                "homepageURL": "https://github.com/snoe/clojure-lsp", 
                "issuesURL": "https://github.com/snoe/clojure-lsp/issues", 
                "docsURL": "https://github.com/snoe/clojure-lsp/blob/master/README.md"
              }
          },
          // ...
      ]
      // ...
  }
  ```

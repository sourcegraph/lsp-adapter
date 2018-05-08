# Experimental Rust language server 

![demonstration of hovers from the Rust language server inside databricks/click](https://cl.ly/383f3V0P1r1u/Screen%20Recording%202018-05-07%20at%2005.58%20PM.gif)

*This GIF was created by browsing https://sourcegraph.com/github.com/databricks/click/-/blob/src/describe.rs with this Rust language server enabled.*

## Introduction

This Dockerfile adds experimental Rust language support for Sourcegraph. 

Thanks to the [rust-lang-nursery/rls](https://github.com/rust-lang-nursery/rls) project for providing the language server that's wrapped by `lsp-adapter` in this image.

## How to add this language server to Sourcegraph

*These steps are adapted from our documentation for [manually adding code intelligence to Sourcegraph server](https://about.sourcegraph.com/docs/code-intelligence/install-manual/).*

1. Run the `sourcegarph/server` Docker image: 

```shell
docker run --publish 7080:7080 --rm --network=lsp --name=sourcegraph --volume ~/.sourcegraph/config:/etc/sourcegraph --volume ~/.sourcegraph/data:/var/opt/sourcegraph sourcegraph/server:2.7.6
```

2. Run the experimental Dockerfile language server:

  ```shell
  docker pull sourcegraph/codeintel-rust

  docker run --rm --network=lsp --name=rust sourcegraph/codeintel-rust
  ```

3. Add the following entry to the `langservers` field in the [site configuration](https://about.sourcegraph.com/docs/config):

  ```js
  {
      // ...
      "langservers": [
          // ...
          {
              "language": "rust",
              "address": "tcp://rust:8080",
              "metadata": {
                "experimental": true,
                "homepageURL": "https://github.com/rust-lang-nursery/rls", 
                "issuesURL": "https://github.com/rust-lang-nursery/rls/issues", 
                "docsURL": "https://github.com/rust-lang-nursery/rls/blob/master/README.md"
              }
          },
          // ...
      ]
      // ...
  }
  ```

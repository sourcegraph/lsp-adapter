# Experimental OCaml language server 

![demonstration of hovers from the OCaml language server inside moby/datakit](https://cl.ly/3l0Z2g1l3D2r/Screen%20Recording%202018-05-07%20at%2005.30%20PM.gif)

*This GIF was created by browsing https://sourcegraph.com/github.com/moby/datakit/-/blob/ci/src/cI_docker.mli with this OCaml language server enabled.*

## Introduction

This Dockerfile adds experimental OCaml language support for Sourcegraph. 

Thanks to the [freebroccolo/ocaml-language-server](https://github.com/freebroccolo/ocaml-language-server) project for providing the language server that's wrapped by `lsp-adapter` in this image.

## How to add this language server to Sourcegraph

*These steps are adapted from our documentation for [manually adding code intelligence to Sourcegraph server](https://about.sourcegraph.com/docs/code-intelligence/install-manual/).*

1. Run the `sourcegarph/server` Docker image: 

```shell
docker run --publish 7080:7080 --rm --network=lsp --name=sourcegraph --volume ~/.sourcegraph/config:/etc/sourcegraph --volume ~/.sourcegraph/data:/var/opt/sourcegraph -v /var/run/docker.sock:/var/run/docker.sock sourcegraph/server:2.7.6
```

2. Run the experimental Dockerfile language server:

  ```shell
  docker pull sourcegraph/codeintel-ocaml

  docker run --rm --network=lsp --name=ocaml sourcegraph/codeintel-ocaml
  ```

3. Add the following entry to the `langservers` field in the [site configuration](https://about.sourcegraph.com/docs/config):

  ```js
  {
      // ...
      "langservers": [
          // ...
          {
              "language": "ocaml",
              "address": "tcp://ocaml:8080",
              "metadata": {
                "experimental": true,
                "homepageURL": "https://github.com/freebroccolo/ocaml-language-server", 
                "issuesURL": "https://github.com/freebroccolo/ocaml-language-server/issues", 
                "docsURL": "https://github.com/freebroccolo/ocaml-language-server/blob/master/README.md"
              }
          },
          // ...
      ]
      // ...
  }
  ```

# Experimental C++ language server 

![demonstration of hovers from the C++ language server inside antirez/redis](https://cl.ly/0f401u080U3S/Screen%20Recording%202018-04-30%20at%2012.21%20pm.gif)

*This GIF was created by browsing https://sourcegraph.com/github.com/antirez/redis/-/blob/src/expire.c with this C++ language server enabled.*

## Introduction

This Dockerfile adds experimental C++ language support for Sourcegraph. 

Thanks to [LLVM](https://clang.llvm.org/extra/clangd.html) (and [Chilledheart/vim-clangd](https://github.com/Chilledheart/vim-clangd) for providing prebuilt binaries) for providing the language server that's wrapped by `lsp-adapter` in this image.

## How to add this language server to Sourcegraph

*These steps are adapted from our documentation for [manually adding code intelligence to Sourcegraph server](https://about.sourcegraph.com/docs/code-intelligence/install-manual/).*

1. Run the `sourcegarph/server` Docker image: 

```shell
docker run --publish 7080:7080 --rm --network=lsp --name=sourcegraph --volume ~/.sourcegraph/config:/etc/sourcegraph --volume ~/.sourcegraph/data:/var/opt/sourcegraph -v /var/run/docker.sock:/var/run/docker.sock sourcegraph/server:2.7.6
```

2. Run the experimental C++ language server:

  ```shell
  docker pull sourcegraph/codeintel-cpp

  docker run --rm --network=lsp --name=cpp sourcegraph/codeintel-cpp
  ```

3. Add the following entry to the `langservers` field in the [site configuration](https://about.sourcegraph.com/docs/config):

  ```js
  {
      // ...
      "langservers": [
          // ...
          {
              "language": "cpp",
              "address": "tcp://cpp:8080",
              "metadata": {
                "experimental": true,
                "homepageURL": "https://clang.llvm.org/extra/clangd.html", 
                "issuesURL": "https://bugs.llvm.org/buglist.cgi?quicksearch=clangd", 
                "docsURL": "https://clang.llvm.org/extra/clangd.html"
              }
          },
          // ...
      ]
      // ...
  }
  ```

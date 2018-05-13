# Experimental Dockerfile language server 

![demonstration of hovers from the Dockerfile language server inside moby/datakit](https://cl.ly/2f3P2t2j2W2i/Screen%20Recording%202018-05-07%20at%2004.55%20PM.gif)

<!-- TODO(@ggilmore @keegancsmith @felixfbecker): Revisit creating this GIF once improved tooltip and syntax highlighting code lands-->

*This GIF was created by browsing https://sourcegraph.com/github.com/moby/datakit/-/blob/ci/self-ci/Dockerfile with this Dockerfile language server enabled.*

## Introduction

This Dockerfile adds experimental Dockerfile language support for Sourcegraph. 

Thanks to the [rcjsuen/dockerfile-language-server-nodejs](https://github.com/rcjsuen/dockerfile-language-server-nodejs) project for providing the language server that's wrapped by `lsp-adapter` in this image.

## How to add this language server to Sourcegraph

*These steps are adapted from our documentation for [manually adding code intelligence to Sourcegraph server](https://about.sourcegraph.com/docs/code-intelligence/install-manual/).*

1. Run the `sourcegarph/server` Docker image: 

```shell
docker run --publish 7080:7080 --rm --network=lsp --name=sourcegraph --volume ~/.sourcegraph/config:/etc/sourcegraph --volume ~/.sourcegraph/data:/var/opt/sourcegraph -v /var/run/docker.sock:/var/run/docker.sock sourcegraph/server:2.7.6
```

2. Run the experimental Dockerfile language server:

  ```shell
  docker pull sourcegraph/codeintel-docker

  docker run --rm --network=lsp --name=dockerfile sourcegraph/codeintel-docker
  ```

3. Add the following entry to the `langservers` field in the [site configuration](https://about.sourcegraph.com/docs/config):

  ```js
  {
      // ...
      "langservers": [
          // ...
          {
              "language": "dockerfile",
              "address": "tcp://dockerfile:8080",
              "metadata": {
                "experimental": true,
                "homepageURL": "https://github.com/rcjsuen/dockerfile-language-server-nodejs", 
                "issuesURL": "https://github.com/rcjsuen/dockerfile-language-server-nodejs/issues", 
                "docsURL": "https://github.com/rcjsuen/dockerfile-language-server-nodejs/blob/master/README.md"
              }
          },
          // ...
      ]
      // ...
  }
  ```

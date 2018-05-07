# Experimental Bash language server 

![demonstration of hovers from the Bash language server inside creationix/nvm](https://cl.ly/271p292i342p/Screen%20Recording%202018-05-07%20at%2009.19%20AM.gif)

*This GIF was created by browsing https://sourcegraph.com/github.com/creationix/nvm/-/blob/nvm.sh with this Bash language server enabled.*

## Introduction

This Dockerfile adds experimental Bash language support for Sourcegraph. 

Thanks to the [mads-hartmann/bash-language-server](https://github.com/mads-hartmann/bash-language-server) project for providing the language server that's wrapped by `lsp-adapter` in this image.

## How to add this language server to Sourcegraph

*These steps are adapted from our documentation for [manually adding code intelligence to Sourcegraph server](https://about.sourcegraph.com/docs/code-intelligence/install-manual/).*

1. Run the `sourcegarph/server` Docker image: 

```shell
docker run \
  --publish 7080:7080 \
  --rm \
  --network=lsp \
  --name=sourcegraph 
  --volume ~/.sourcegraph/config:/etc/sourcegraph \
  --volume ~/.sourcegraph/data:/var/opt/sourcegraph \
  sourcegraph/server:2.7.6
```

2. Run the experimental Bash language server:

  ```shell
  docker pull sourcegraph/codeintel-bash

  docker run --rm --network=lsp --name=bash sourcegraph/codeintel-bash
  ```

3. Add the following entry to the `langservers` field in the [site configuration](https://about.sourcegraph.com/docs/config):

  ```json
  {
      // ...
      "langservers": [
          // ...
          {
              "language": "bash",
              "address": "tcp://bash:8080",
              "metadata": {
                "experimental": true,
                "homepageURL": "https://github.com/mads-hartmann/bash-language-server", 
                "issuesURL": "https://github.com/mads-hartmann/bash-language-server/issues", 
                "docsURL": "https://github.com/mads-hartmann/bash-language-server/blob/master/README.md"
              }
          },
          // ...
      ]
      // ...
  }
  ```

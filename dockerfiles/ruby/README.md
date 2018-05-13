# Experimental Ruby language server 

![demonstration of hovers from the Ruby language server inside tabler/tabler](https://cl.ly/3l0Z2g1l3D2r/Screen%20Recording%202018-05-07%20at%2005.30%20PM.gif)

*This GIF was created by browsing https://sourcegraph.com/github.com/tabler/tabler/-/blob/src/_plugins/jekyll-tabler.rb with this Ruby language server enabled.*

## Introduction

This Dockerfile adds experimental Ruby language support for Sourcegraph. 

Thanks to the [castwide/solargraph](https://github.com/castwide/solargraph) project for providing the language server that's wrapped by `lsp-adapter` in this image.

## How to add this language server to Sourcegraph

*These steps are adapted from our documentation for [manually adding code intelligence to Sourcegraph server](https://about.sourcegraph.com/docs/code-intelligence/install-manual/).*

1. Run the `sourcegarph/server` Docker image: 

```shell
docker run --publish 7080:7080 --rm --network=lsp --name=sourcegraph --volume ~/.sourcegraph/config:/etc/sourcegraph --volume ~/.sourcegraph/data:/var/opt/sourcegraph -v /var/run/docker.sock:/var/run/docker.sock sourcegraph/server:2.7.6
```

2. Run the experimental Dockerfile language server:

  ```shell
  docker pull sourcegraph/codeintel-ruby

  docker run --rm --network=lsp --name=ruby sourcegraph/codeintel-ruby
  ```

3. Add the following entry to the `langservers` field in the [site configuration](https://about.sourcegraph.com/docs/config):

  ```js
  {
      // ...
      "langservers": [
          // ...
          {
              "language": "ruby",
              "address": "tcp://ruby:8080",
              "metadata": {
                "experimental": true,
                "homepageURL": "https://github.com/castwide/solargraph", 
                "issuesURL": "https://github.com/castwide/solargraph/issues", 
                "docsURL": "https://github.com/castwide/solargraph/blob/master/README.md"
              }
          },
          // ...
      ]
      // ...
  }
  ```

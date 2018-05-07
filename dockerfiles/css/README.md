# Experimental CSS language server 

![demonstration of hovers from the CSS language server inside tabler/tabler](https://cl.ly/3I2S3x2Z3K44/Screen%20Recording%202018-05-07%20at%2004.36%20PM.gif)

*This GIF was created by browsing https://sourcegraph.com/github.com/tabler/tabler/-/blob/dist/assets/plugins/prismjs/plugin.css with this CSS language server enabled.*

## Introduction

This Dockerfile adds experimental CSS language support for Sourcegraph. 

Thanks to the [Microsoft/vscode](https://github.com/Microsoft/vscode/) project (and [vscode-langservers/vscode-css-languageserver-bin](https://github.com/vscode-langservers/vscode-css-languageserver-bin) for providing prebuilt binaries) for providing the language server that's wrapped by `lsp-adapter` in this image.

## How to add this language server to Sourcegraph

*These steps are adapted from our documentation for [manually adding code intelligence to Sourcegraph server](https://about.sourcegraph.com/docs/code-intelligence/install-manual/).*

1. Run the `sourcegarph/server` Docker image: 

```shell
docker run --publish 7080:7080 --rm --network=lsp --name=sourcegraph --volume ~/.sourcegraph/config:/etc/sourcegraph --volume ~/.sourcegraph/data:/var/opt/sourcegraph sourcegraph/server:2.7.6
```

2. Run the experimental CSS language server:

  ```shell
  docker pull sourcegraph/codeintel-css

  docker run --rm --network=lsp --name=css sourcegraph/codeintel-css
  ```

3. Add the following entry to the `langservers` field in the [site configuration](https://about.sourcegraph.com/docs/config):

  ```json
  {
      // ...
      "langservers": [
          // ...
          {
              "language": "css",
              "address": "tcp://css:8080",
              "metadata": {
                "experimental": true,
                "homepageURL": "https://github.com/Microsoft/vscode", 
                "issuesURL": "https://github.com/Microsoft/vscode/issues", 
                "docsURL": "https://github.com/Microsoft/vscode/blob/master/README.md"
              }
          },
          // ...
      ]
      // ...
  }
  ```

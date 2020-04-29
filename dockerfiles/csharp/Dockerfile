FROM golang:1.14-alpine AS lsp-adapter
WORKDIR /go/src/github.com/sourcegraph/lsp-adapter
COPY . .
RUN CGO_ENABLED=0 GOBIN=/usr/local/bin go install github.com/sourcegraph/lsp-adapter

# .NET Core 2.0.0 is not available on Alpine yet, but it is actively being
# worked on: https://github.com/dotnet/core/issues/1076
FROM ubuntu:16.04

RUN apt-get update \
  && apt-get install --no-install-recommends -y ca-certificates gnupg curl jq \
  # Ubuntu 16.04's default nodejs package v4.2.6 is too old to run the
  # omnisharp-client npm package.
  && curl -sL https://deb.nodesource.com/setup_10.x | bash - \
  && curl https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > /etc/apt/trusted.gpg.d/microsoft.gpg \
  && echo "deb [arch=amd64] https://packages.microsoft.com/repos/microsoft-ubuntu-xenial-prod xenial main" > /etc/apt/sources.list.d/dotnetdev.list \
  && echo "deb https://download.mono-project.com/repo/ubuntu stable-xenial main" > /etc/apt/sources.list.d/mono-official-stable.list \
  && apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 3FA7E0328081BFF6A14DA29AA6A19B38D3D831EF \
  && apt-get update && apt-get install --no-install-recommends -y apt-transport-https \
  # OmniSharp runs on Mono and requires the .NET SDK 2.0.0.
  && apt-get update && apt-get install --no-install-recommends -y mono-complete dotnet-sdk-2.0.0 nodejs \
  # Removes ~1 GB of unused files.
  && rm -rf /usr/share/dotnet/sdk/NuGetFallbackFolder \
  && rm -rf /var/lib/apt/lists/* \
  # omnisharp-client translates between LSP and OmniSharp's custom protocol.
  && npm install -g omnisharp-client@7.2.3

# Bake OmniSharp into this image (and write a .version file) so that
# omnisharp-client doesn't have to download it upon initialization.
ENV BASE=/usr/lib/node_modules/omnisharp-client
RUN mkdir "$BASE/omnisharp-linux-x64"
RUN VERSION="$(jq -r ".[\"omnisharp-roslyn\"]" < "$BASE/package.json")" \
  && curl -L "https://github.com/OmniSharp/omnisharp-roslyn/releases/download/$VERSION/omnisharp-linux-x64.tar.gz" -o "$BASE/omnisharp-linux-x64.tar.gz" \
  && echo "$VERSION" > "$BASE/omnisharp-linux-x64/.version" \
  && tar xvf "$BASE/omnisharp-linux-x64.tar.gz" -C "$BASE/omnisharp-linux-x64"

# Use tini as entrypoint to correctly handle signals and zombie processes.
ENV TINI_VERSION v0.18.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]

EXPOSE 8080
COPY --from=lsp-adapter /usr/local/bin/lsp-adapter /usr/local/bin/
CMD ["lsp-adapter", "-trace", "-proxyAddress=0.0.0.0:8080", "nodejs", "/usr/lib/node_modules/omnisharp-client/languageserver/server.js"]

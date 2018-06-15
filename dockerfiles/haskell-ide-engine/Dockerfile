# This Dockerfile compiles haskell-ide-engine to /root/.local/bin/hie.
#
# The compilation process requires more than 3 GB of RAM, so make sure you set
# your Docker preferences to ~8 GB of RAM before building this image.

FROM ubuntu:16.04

ARG STACK_VERSION=1.7.1

# --no-install-recommends actually breaks tar.
# hadolint ignore=DL3015
RUN apt-get update && apt-get install -y \
  git \
  wget \
  libtinfo-dev \
  build-essential \
  libgmp3-dev \
  zlib1g-dev \
  && rm -rf /var/lib/apt/lists/*

RUN wget -qO- "https://github.com/commercialhaskell/stack/releases/download/v$STACK_VERSION/stack-$STACK_VERSION-linux-x86_64.tar.gz" | tar xz --wildcards --strip-components=1 -C /usr/local/bin '*/stack'

# hadolint ignore=DL3003
RUN git clone https://github.com/haskell/haskell-ide-engine --recursive /tmp/haskell-ide-engine \
  && cd /tmp/haskell-ide-engine \
  && git checkout 562ac94d245e7b6ffa380eae4b02a832c397cfbb \
  # Avoid invalidating the layers when new commits are added.
  && find /tmp/haskell-ide-engine -name .git -print0 | xargs -0 rm -rf

WORKDIR /tmp/haskell-ide-engine

RUN stack install
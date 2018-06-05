FROM golang:1.10-alpine AS lsp-adapter
WORKDIR /go/src/github.com/sourcegraph/lsp-adapter
COPY . .
RUN CGO_ENABLED=0 GOBIN=/usr/local/bin go install github.com/sourcegraph/lsp-adapter

FROM mongo:3.7.9

RUN apt-get update \
  && apt-get install --no-install-recommends -y curl=7.47.0-1ubuntu2.8 \
  && rm -rf /var/lib/apt/lists/*
RUN curl -sL https://deb.nodesource.com/setup_8.x | bash -
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update \
  && apt-get install --no-install-recommends -y \
    nodejs=8.11.2-1nodesource1 \
    man-db=2.7.5-1 \
    netcat=1.10-41 \
    python=2.7.12-1~16.04 \
    python-pip=8.1.1-2ubuntu0.4 \
    git=1:2.7.4-0ubuntu1.3 \
    yarn=1.7.0-1 \
  && rm -rf /var/lib/apt/lists/*

# Downgrade pip because explainshell depends on this API.
RUN pip install --upgrade pip==9.0.3

# This should be changed to the upstream https://github.com/idank/explainshell
# once https://github.com/idank/explainshell/pull/125 is merged in.
#
# Ignore the warning about using WORKDIR instead of cd for convenience.
# hadolint ignore=DL3003
RUN git clone https://github.com/chrismwendt/explainshell \
  && cd explainshell \
  && git checkout 775a6097ab19ec9ee320919f7f29f7041513125e \
  && rm -rf explainshell/.git \
  && pip install setuptools==39.2.0 \
  && pip install -r requirements.txt \
  # Prepare for baking the man page data into this image by pointing mongo at a
  # db directory other than /data/db, which is read-only in the upstream base
  # mongo image. See https://docs.docker.com/engine/reference/builder/#volume
  # "If any build steps change the data within the volume after it has been
  # declared, those changes will be discarded."
  && mkdir -p /data/db2 \
  && echo "dbpath = /data/db2" > /etc/mongodb.conf \
  && chown -R mongodb:mongodb /data/db2

# Use tini as entrypoint to correctly handle signals and zombie processes.
ENV TINI_VERSION v0.18.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]

RUN apt-get update && apt-get install --no-install-recommends -y \
	make=4.1-6 \
	gcc=4:5.3.1-1ubuntu1 \
	g++=4:5.3.1-1ubuntu1 \
  && rm -rf /var/lib/apt/lists/*

# This should be changed to the upstream
# https://github.com/mads-hartmann/bash-language-server once
# https://github.com/mads-hartmann/bash-language-server/pull/45 is merged in.
#
# Ignore the warning about using WORKDIR instead of cd for convenience.
# hadolint ignore=DL3003
RUN git clone https://github.com/chrismwendt/bash-language-server \
  && cd bash-language-server \
  && git checkout explainshell \
  && make build \
  && make install \
  && yarn global add file:"$(pwd)"/server

COPY . .

# Very slow (~30 minutes) - initializes the mongo database with thousands of man
# pages and bakes the db into the image (ends up being ~2 GB).
#
# Ignore the warning about using WORKDIR instead of cd for convenience.
# hadolint ignore=DL3003
RUN cd explainshell && ../dockerfiles/bash/populate-explainshell-db.sh

# https://stackoverflow.com/a/33601894
# This might not be necessary because this db is read-only
# TODO(chris): Try removing this line and seeing if it still works
# VOLUME /data/db2

COPY --from=lsp-adapter /usr/local/bin/lsp-adapter /usr/local/bin
EXPOSE 8080

CMD ["./dockerfiles/bash/run-bash-language-server.sh"]

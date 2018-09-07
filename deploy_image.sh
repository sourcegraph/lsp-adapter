#!/bin/bash
set -e

DOCKER_REPOSITORY="sourcegraph/codeintel-$LANGUAGE"
DOCKERFILE_PATH="./dockerfiles/$LANGUAGE/Dockerfile"
VERSION=$(printf "%05d" $TRAVIS_JOB_NUMBER)_$(date +%Y-%m-%d)_$(git rev-parse --short HEAD)

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker build \
	-t $DOCKER_REPOSITORY:insiders \
	-t $DOCKER_REPOSITORY:latest \
	-t $DOCKER_REPOSITORY:$VERSION \
	-f $DOCKERFILE_PATH .
docker push $DOCKER_REPOSITORY
docker push $DOCKER_REPOSITORY:$VERSION
docker push $DOCKER_REPOSITORY:insiders

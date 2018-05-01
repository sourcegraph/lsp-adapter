#!/bin/bash
set -e

DOCKER_REPOSITORY="sourcegraph/codeintel-$LANGUAGE"
DOCKERFILE_PATH="./dockerfiles/$LANGUAGE/Dockerfile"

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker build \
	-t ${DOCKER_REPOSITORY}:latest \
	-t ${DOCKER_REPOSITORY}:${TRAVIS_COMMIT} \
	-t ${DOCKER_REPOSITORY}:${TRAVIS_JOB_NUMBER} \
	-t ${DOCKER_REPOSITORY}:$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
	-f ${DOCKERFILE_PATH} .
docker push $DOCKER_REPOSITORY
docker push ${DOCKER_REPOSITORY}:${TRAVIS_COMMIT}
docker push ${DOCKER_REPOSITORY}:${TRAVIS_JOB_NUMBER}
docker push ${DOCKER_REPOSITORY}:$(date -u +"%Y-%m-%dT%H:%M:%SZ")

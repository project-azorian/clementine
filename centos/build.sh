#!/bin/bash
set -ex

: "${DISTRO:="centos"}"
: "${DOCKER_REGISTRY:="docker.io"}"
: "${DOCKER_REPO:="azorian"}"
: "${DOCKER_CLEMENTINE_PREFIX:="clementine"}"
: "${DOCKER_BUILD_IMAGE:="${DOCKER_CLEMENTINE_PREFIX}-${DISTRO}-builder"}"
: "${DOCKER_IMAGE:="${DOCKER_CLEMENTINE_PREFIX}-${DISTRO}"}"
: "${DOCKER_TAG:="latest"}"
: "${DOCKER_BUILD_IMAGE_FQN:="${DOCKER_REGISTRY}/${DOCKER_REPO}/${DOCKER_BUILD_IMAGE}:${DOCKER_TAG}"}"
: "${DOCKER_IMAGE_FQN:="${DOCKER_REGISTRY}/${DOCKER_REPO}/${DOCKER_IMAGE}:${DOCKER_TAG}"}"

docker build \
  --network="host" \
  --tag="${DOCKER_BUILD_IMAGE_FQN}" \
  --file="./builder/Dockerfile" \
  .

docker run \
  --privileged \
  --net="host" \
  --name="${DOCKER_REPO}-${DOCKER_BUILD_IMAGE}-${DOCKER_TAG}" \
  ${DOCKER_BUILD_IMAGE_FQN} \
    clementine-build

docker commit "${DOCKER_REPO}-${DOCKER_BUILD_IMAGE}-${DOCKER_TAG}" "${DOCKER_IMAGE_FQN}"
docker rm -f "${DOCKER_REPO}-${DOCKER_BUILD_IMAGE}-${DOCKER_TAG}"
docker push "${DOCKER_IMAGE_FQN}"

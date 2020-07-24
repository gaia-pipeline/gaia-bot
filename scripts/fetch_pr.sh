#!/bin/bash

set -e
set -u
set -o pipefail

pr="<pr_replace>"
branch="<branch_replace>"
folder=$(mktemp -d -t fetch-XXXXXXXXXX)
tag="<tag_replace>"
docker_token="<docker_token_replace>"
docker_username="<docker_username_replace>"

function main() {
  if [[ -z "${pr}" ]]; then
    echo "PR_NUMBER number is empty. Please set.";
    exit 1
  fi

  if [[ -z "${branch}" ]]; then
    echo "BRANCH is empty. Please set.";
    exit 1
  fi

  if [[ -z "${folder}" ]]; then
    echo "FOLDER is empty. Please set.";
    exit 1
  fi

  if [[ -z "${docker_token}" ]]; then
    echo "DOCKER_TOKEN is empty. Please set.";
    exit 1
  fi

  if [[ -z "${docker_username}" ]]; then
    echo "DOCKER_USERNAME is empty. Please set.";
    exit 1
  fi

  if [[ -z "${tag}" ]]; then
    echo "TAG is empty. Please set.";
    exit 1
  fi

  # Check out PR
  cd "${folder}"
  mkdir -p gaia
  cd gaia
  git clone https://github.com/gaia-pipeline/gaia.git
  cd gaia
  git fetch origin pull/"${pr}"/head:"${branch}"
  git checkout "${branch}"

  # Make a release
  make download
  make release


  # build docker image
  docker build -t "${tag}" . -f ./docker/Dockerfile

  # push image to gaia test repo
  echo "${docker_token}" | docker login --username "${docker_username}" --password-stdin
  docker push "${tag}"
}

trap cleanup EXIT

function cleanup {
  # Docker prune system and earase gaia image after push.
  docker system prune --force
  docker rmi "${tag}" --force || true
}

# Run the script
main
#!/bin/bash

set -e
set -u
set -o pipefail

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

# Parse the opt values.
while getopts ':p:t:b:f:o:u:' OPTION; do
  case "$OPTION" in
    p)
      pr="$OPTARG"
      ;;

    t)
      tag="$OPTARG"
      ;;

    b)
      branch="$OPTARG"
      ;;

    f)
      folder="$OPTARG"
      ;;

    o)
      docker_token="$OPTARG"
      ;;

    u)
      docker_username="$OPTARG"
      ;;
    ?)
      echo "script usage: $(basename "${0}") [-p pr] [-t tag] [-b branch] [-f folder] [-o docker_token] [-u docker_username]" >&2
      exit 1
      ;;
  esac
done
shift "$(($OPTIND -1))"

# Run the script
main
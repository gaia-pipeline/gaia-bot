#!/bin/bash

set -e
set -u
set -o pipefail

function main() {
  if [[ -z "${tag}" ]]; then
    echo "Tag is empty. Please set.";
    exit 1
  fi

  if [[ -z "${folder}" ]]; then
    echo "Folder is empty. Please set.";
    exit 1
  fi

  if [[ -z "${git_token}" ]]; then
    echo "GIT_TOKEN is empty. Please set.";
    exit 1
  fi

  if [[ -z "${git_username}" ]]; then
    echo "GIT_USERNAME is empty. Please set.";
    exit 1
  fi

  # checkout infra code
  cd "${folder}"
  mkdir -p infra
  cd infra
  git clone https://"${git_username}":"${git_token}"@"${repo}" infra
  cd infra
  sed -i "s/image:.*/image: ${tag}/g" workloads/gaia_deployment.yaml
  git commit -am 'Updated tag for gaia deployment'
  git push origin master
}

# Parse the opt values.
while getopts ':r:t:f:o:u:' OPTION; do
  case "$OPTION" in
    t)
      tag="$OPTARG"
      ;;

    r)
      repo="$OPTARG"
      ;;

    f)
      folder="$OPTARG"
      ;;

    o)
      git_token="$OPTARG"
      ;;

    u)
      git_username="$OPTARG"
      ;;
    ?)
      echo "script usage: $(basename "${0}") [-r repo] [-t tag] [-f folder] [-o git_token] [-u git_username]" >&2
      exit 1
      ;;
  esac
done
shift "$(($OPTIND -1))"

# Run the script
main
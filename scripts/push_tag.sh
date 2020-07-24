#!/bin/bash

set -e
set -u
set -o pipefail

repo="<repo_replace>"
folder=$(mktemp -d -t push-XXXXXXXXXX)
tag="<tag_replace>"
git_token="<git_token_replace>"
git_username="<git_username_replace>"

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

# Run the script
main
#!/bin/bash

set -e

tag="${TAG}"
folder="${FOLDER}"
repo="${REPO}"

if [[ -z "${tag}" ]]; then
  echo "Tag is empty. Please set.";
  exit 1
fi

if [[ -z "${folder}" ]]; then
  echo "Folder is empty. Please set.";
  exit 1
fi

# checkout infra code
cd "${folder}"
git clone "${repo}"



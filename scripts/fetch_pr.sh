#!/bin/bash

set -e

pr="${PR_NUMBER}"
repo="${REPOSITORY}"
branch="${BRANCH}"
folder="${FOLDER}"
docker_token="${DOCKER_TOKEN}"
docker_username="${DOCKER_USERNAME}"

if [[ -z "${pr}" ]]; then
  echo "PR number is empty. Please set.";
  exit 1
fi

if [[ -z "${repo}" ]]; then
  echo "REPOSITORY empty. Please set.";
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

# Check out PR
cd "${folder}"
mkdir -p gaia
cd gaia
git clone https://github.com/gaia-pipeline/gaia.git
cd gaia
git fetch origin pull/"${pr}"/head:"${branch}"
git checkout "${branch}"

# Make release
wget -qO- https://raw.githubusercontent.com/creationix/nvm/v0.33.11/install.sh | bash
echo 'export NVM_DIR=$HOME/.nvm' >> ~/.bash_profile
touch $HOME/.nvmrc
echo 'source $NVM_DIR/nvm.sh' >> ~/.bash_profile
source ~/.bash_profile
nvm install v12.6.0
npm cache clean --force
make download
make release

# build docker image
tag="gaiapipeline/testing:${branch}-${pr}"
docker build -t "${tag}" ./docker -f ./docker/Dockerfile

# push image to gaia test repo
echo "${docker_token}" | docker login --username "${docker_username}" --password-stdin
docker push "${tag}"
# This output will be used by the gaia-bot to signal flux to deploy this image once it's done pushing.
echo "${tag}"
#!/bin/bash

set -e

pr="${PR_NUMBER}"
tag="${TAG}"
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

# Make release
service docker start
wget -qO- https://raw.githubusercontent.com/creationix/nvm/v0.33.11/install.sh | bash
echo 'export NVM_DIR=$HOME/.nvm' >> ~/.bash_profile
touch $HOME/.nvmrc
echo 'source $NVM_DIR/nvm.sh' >> ~/.bash_profile
source ~/.bash_profile
nvm install v12.6.0 || true
npm cache clean --force || true
make download
make release

trap cleanup EXIT

function cleanup {
  # Docker prune system and earase gaia image after push.
  docker system prune --force
  docker rmi "${tag}" --force
}

# build docker image
docker build -t "${tag}" . -f ./docker/Dockerfile

# push image to gaia test repo
echo "${docker_token}" | docker login --username "${docker_username}" --password-stdin
docker push "${tag}"

#!/bin/bash

set -e

pr="${PR_NUMBER}"
repo="${REPOSITORY}"
branch="${BRANCH}"
folder="${FOLDER}"
docker_token="${DOCKER_TOKEN}"

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

# Check out PR
cd "${folder}"


# Compile frontend
wget -qO- https://raw.githubusercontent.com/creationix/nvm/v0.33.11/install.sh | bash
echo 'export NVM_DIR=$HOME/.nvm' >> ~/.bash_profile
touch $HOME/.nvmrc
echo 'source $NVM_DIR/nvm.sh' >> ~/.bash_profile
source ~/.bash_profile
nvm install v12.6.0
npm cache clean --force
make download
make release

# Compile backend
make compile_backend

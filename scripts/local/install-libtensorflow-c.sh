#!/bin/sh
set -x
set -e

clear
echo

## tensorflow vars
TF_VERSION=${TF_VERSION:-"1.3.0"} 	# Change for the version you need
TF_TYPE=${TF_TYPE:-"cpu"}			# Change to "gpu" for GPU support

TF_VCS_REMOTE_URI=${TF_VCS_REMOTE_URI:-"github.com/tensorflow/tensorflow"}
TF_VCS_CLONE_DIR=${TF_VCS_CLONE_DIR:-"$GOPATH/$TF_VCS_REMOTE_URI"}
TF_VCS_CLONE_DEPTH=${TF_VCS_CLONE_DEPTH:-"1"}


## local vars
TARGET_DIRECTORY='/usr/local'
HOME_BIN=$HOME/bin

## compilers vars
GCC_VERSION=$(gcc --version)

## path(s) for executables
export PATH="$PATH:${HOME_BIN}"

## install bazel via homebrew
brew install bazel
# brew upgrade bazel
brew install python3

## download tensorflow source code from github.com
# rm -fR ${TF_VCS_CLONE_DIR}
if [ ! -d "${TF_VCS_CLONE_DIR}" ]; then
	git clone --recursive --depth=${TF_VCS_CLONE_DEPTH} -b v${TF_VERSION} https://${TF_VCS_REMOTE_URI} ${TF_VCS_CLONE_DIR}
fi
cd ${TF_VCS_CLONE_DIR}
pwd

## pre-requisites
python3 -m ensurepip
sudo pip3 install --no-cache --upgrade cython six numpy wheel virtualenv tensorflow

## build from source
#  /usr/local/bin/python3
./configure
bazel build --config=opt //tensorflow/tools/pip_package:build_pip_package
bazel-bin/tensorflow/tools/pip_package/build_pip_package /tmp/tensorflow_pkg

# virtualenv --system-site-packages -p python3

exit 1

## install libtensorflow c-library for TF's golang package setup
curl -L \
"https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-${TF_TYPE}-$(go env GOOS)-x86_64-${TF_VERSION}.tar.gz" |
sudo tar -C $TARGET_DIRECTORY -xz

sudo ldconfig
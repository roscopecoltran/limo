#!/bin/bash

# git subtree add templates/tensorflow/application --prefix https://github.com/tobegit3hub/tensorflow_template_application master --squash

set -ex
# install gRPC and protoc plugin for Go, see http://www.grpc.io/docs/quickstart/go.html#generate-grpc-code
mkdir tensorflow tensorflow_serving
protoc -I model/generate/ model/generate/tensorflow/*.proto --go_out=plugins=grpc:tf_server
protoc -I model/generate/ model/generate/tensorflow/core/framework/* --go_out=plugins=grpc:.
#!/bin/sh
set -e
go-bindata templates licenses
go build -ldflags="-w -s"
go fmt

rm -rf ./workbench*

# application test
mkdir workbench1
pushd workbench1
../qtpm init app
../qtpm build
popd
# library test
mkdir workbench2
pushd workbench2
../qtpm init lib
../qtpm build
popd


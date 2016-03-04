#!/bin/sh
go-bindata templates licenses
go build -ldflags="-w -s"
rm -rf ./workbench
mkdir workbench
cd workbench
../qtpm init app HelloWorld mit
../qtpm build

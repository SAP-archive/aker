#!/bin/bash

set -e

mkdir -p $GOPATH/src

echo "Moving project to GOPATH..."
prefix_path=$GOPATH/src/github.infra.hana.ondemand.com/I061150
mkdir -p $prefix_path
cp -r aker $prefix_path
cd $prefix_path/aker

echo "Fetching test tools..."
go get github.com/onsi/ginkgo/ginkgo

echo "Running tests..."
ginkgo -r

#!/usr/bin/env bash

fail() {
  echo "unit tests failed"
  exit 1
}

FILES=$(go list "$1" | grep -v /vendor/) || fail
go test -tags=unit -timeout 30s -short -v ${FILES} || fail
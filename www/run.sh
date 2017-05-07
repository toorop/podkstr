#!/bin/sh

set -e;

go build -o dist/server

gulp build

dist/server
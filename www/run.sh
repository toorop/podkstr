#!/bin/sh

set -e;

go build -o dist/server

gulp minify

dist/server
#!/bin/sh

set -e;

go build -o dist/server
gulp build

rsync -arvz --exclude=config.base.yml --exclude=config.yml dist/* root@podkstr.com:/var/www/podkstr.com/

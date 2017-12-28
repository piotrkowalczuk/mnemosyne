#!/usr/bin/env bash

touch c.out
rm c.out
set -e

bash <(curl -s https://codecov.io/bash)
gocov convert coverage.out > c.out
./cc-test-reporter after-build

docker login -u $DOCKER_USER -p $DOCKER_PASSWORD
if [ ! -z "$TRAVIS_TAG" ] && [ "$TRAVIS_PULL_REQUEST" == "false" ]; then
    export VCS_REF=$TRAVIS_TAG
    make publish
fi
if [ $TRAVIS_BRANCH == 'master' ] && [ "$TRAVIS_PULL_REQUEST" == "false" ]; then
    export VERSION=latest
    export VCS_REF=master
    make publish
fi
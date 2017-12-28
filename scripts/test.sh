#!/usr/bin/env bash

touch coverage.out
rm coverage.out
set -e
echo "mode: atomic" > coverage.out

for d in $(go list ./... | grep -v /vendor/ | grep -v /mnemosynerpc| grep -v /mnemosynetest); do
	go test -race -coverprofile=profile.out -covermode=atomic $d
	if [ -f profile.out ]; then
		tail -n +2 profile.out >> coverage.out
		rm profile.out
	fi
done


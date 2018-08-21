#!/usr/bin/env bash

: ${TEST_RESULTS:=.}

touch ${TEST_RESULTS}/coverage.out
rm ${TEST_RESULTS}/coverage.out

set -e
echo "mode: atomic" > ${TEST_RESULTS}/coverage.out

for d in $(go list ./... | grep -v /vendor/ | grep -v /mnemosynerpc| grep -v /mnemosynetest); do
	go test -race -coverprofile=profile.out -covermode=atomic -v ${d}
	if [ -f profile.out ]; then
		tail -n +2 profile.out >> ${TEST_RESULTS}/coverage.out
		rm profile.out
	fi
done


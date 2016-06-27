#!/usr/bin/env bash

set -e
echo "" > coverage.txt

ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --race --compilers=2 -v
for f in $(find . -maxdepth 10 -type f -name \*.coverprofile); do
  if [ -f "$f" ]; then
    cat $f >> coverage.txt
    rm $f
  fi
done

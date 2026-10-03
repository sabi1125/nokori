#!/bin/sh
# Dev only: rebuilds public/openapi.yaml from spec/ whenever any spec
# file changes. Scalar can't follow $refs into other files, so it has to
# be served one bundled file. Polls instead of using inotify so it also
# works on bind mounts where file events don't propagate.
set -u

last=""
while true; do
  now=$(find spec -type f -name '*.yaml' -exec md5sum {} + | sort | md5sum)
  if [ "$now" != "$last" ]; then
    if redocly bundle spec/openapi.yaml -o public/openapi.yaml; then
      echo "bundled at $(date +%T)"
    else
      echo "bundle FAILED at $(date +%T)"
    fi
    last=$now
  fi
  sleep 1
done

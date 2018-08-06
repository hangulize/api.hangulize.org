#!/usr/bin/env bash
set -euo pipefail

# Activate the 'gcloud' command.
source ./google-cloud-sdk/path.bash.inc

# Sign in with the service account key.
# shellcheck disable=SC2154
openssl aes-256-cbc \
  -K "$encrypted_5b7fb06a6635_key" -iv "$encrypted_5b7fb06a6635_iv" \
  -in .service-account.json.enc -out .service-account.json -d
gcloud auth activate-service-account --key-file .service-account.json

# Deploy 'hangulize-api'.
gcloud config set core/project hangulize-api
gcloud app deploy ./dispatch.yaml ./*/app.yaml -q --version="$TRAVIS_COMMIT"

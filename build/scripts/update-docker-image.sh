#!/usr/bin/env bash
#
# This script builds and pushes the pokesay image to ghcr.io.
# This image is primarily used for testing by the GitHub Actions workflow.

set -euo pipefail

# login to ghcr.io
echo $GITHUB_API_TOKEN | docker login ghcr.io -u tmck-code --password-stdin

# build and push the image
docker build -f Dockerfile -t ghcr.io/tmck-code/pokesay:latest ..
docker push ghcr.io/tmck-code/pokesay:latest

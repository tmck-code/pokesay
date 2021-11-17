#!/bin/bash

# Run pokesay and get the timings only

fortune | "${1:-./pokesay}" 2>&1 > /dev/null | jq .

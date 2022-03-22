#!/bin/bash

# Run pokesay and get the timings only

echo w | "${1:-./pokesay}" 2>&1 > /dev/null | jq .

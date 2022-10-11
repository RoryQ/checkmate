#!/bin/bash

set -xe

BIN_VERSION=""
if [ "$INPUT_BIN_VERSION" != "latest" ]; then
  BIN_VERSION="/tags/${$INPUT_BIN_VERSION}"
fi

PLATFORM="$(uname -s)"
ARCH="$(uname -m)"

if [PLATFORM == "linux"] && [ARCH == ""]

wget -O - -q "$(wget -q https://api.github.com/repos/aquasecurity/tfsec/releases${BIN_VERSION} -O - | grep -m 1 -o -E "https://.+?tfsec-linux-amd64" | head -n1)" > tfsec-linux-amd64
wget -O - -q "$(wget -q https://api.github.com/repos/aquasecurity/tfsec/releases${BIN_VERSION} -O - | grep -m 1 -o -E "https://.+?tfsec_checksums.txt" | head -n1)" > tfsec.checksums



tfsec --out=${TFSEC_OUT_OPTION} --format=${TFSEC_FORMAT_OPTION} --soft-fail ${TFSEC_ARGS_OPTION} "${INPUT_WORKING_DIRECTORY}"

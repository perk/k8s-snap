#!/usr/bin/env bash

DIR=`realpath $(dirname "${0}")`

# Initialize node for e2e tests
"${DIR}/connect-interfaces.sh"
"${DIR}/network-requirements.sh"

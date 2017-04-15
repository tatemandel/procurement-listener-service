#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

go build procurementlistenerservice
./procurementlistenerservice --metadataFile ${DIR}/../sample/metadata.json

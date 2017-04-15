#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

env GOOS=linux GOARCH=amd64 go build procurementlistenerservice

if [ $? -eq 0 ]; then
	gsutil cp ./procurementlistenerservice gs://cloud-commerce-procurement
	gsutil cp ${DIR}/../sample/metadata.json gs://cloud-commerce-procurement
fi

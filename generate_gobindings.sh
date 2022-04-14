#!/bin/sh

rm -rf openapi/
openapi-generator generate -g go -i pinning-service.yaml -o openapi --global-property skipFormModel=false --enable-post-process-file
rm openapi/go.mod openapi/go.sum

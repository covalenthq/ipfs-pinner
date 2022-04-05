#!/bin/sh

rm -rf openapi/
openapi-generator generate -g go -i pinning-service.yaml -o openapi --global-property skipFormModel=false
rm openapi/go.mod openapi/go.sum

#!/bin/sh
set -ex
DYNAMODB_ENDPOINT=$1
aws dynamodb create-table --endpoint-url ${DYNAMODB_ENDPOINT} \
  --billing-mode PAY_PER_REQUEST \
  --table-name lifelog-metrics \
  --attribute-definitions \
    AttributeName=partition_key,AttributeType=S \
    AttributeName=sort_key,AttributeType=S \
  --key-schema \
    AttributeName=partition_key,KeyType=HASH \
    AttributeName=sort_key,KeyType=RANGE
aws dynamodb create-table --endpoint-url ${DYNAMODB_ENDPOINT} \
  --billing-mode PAY_PER_REQUEST \
  --table-name lifelog \
  --attribute-definitions \
    AttributeName=partition_key,AttributeType=S \
    AttributeName=sort_key,AttributeType=S \
  --key-schema \
    AttributeName=partition_key,KeyType=HASH \
    AttributeName=sort_key,KeyType=RANGE

#!/bin/bash

aws dynamodb create-table \
    --region us-east-1 \
    --table-name payments \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 \
    --endpoint-url=http://localhost:4566
aws dynamodb put-item \
    --region us-east-1 \
    --table-name payments \
    --item '{
        "id": {"S": "payment1"},
        "amount": {"N": "100"},
        "currency": {"S": "USD"},
        "payment_method": {"S": "credit_card"},
        "status": {"S": "completed"}
    }' \
    --endpoint-url=http://localhost:4566
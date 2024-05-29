#!/bin/bash

awslocal sqs create-queue --queue-name maked-order --region us-east-1
awslocal sqs create-queue --queue-name confirmed-payment --region us-east-1
awslocal sqs create-queue --queue-name notification --region us-east-1

awslocal dynamodb create-table \
    --table-name order \
    --key-schema AttributeName=ID,KeyType=HASH \
    --attribute-definitions AttributeName=ID,AttributeType=S \
    --billing-mode PAY_PER_REQUEST \
    --region us-east-1

awslocal dynamodb create-table \
    --table-name payment \
    --key-schema AttributeName=ID,KeyType=HASH \
    --attribute-definitions AttributeName=ID,AttributeType=S \
    --billing-mode PAY_PER_REQUEST \
    --region us-east-1

#!/bin/bash

export $(grep -v '^#' .env | xargs)

curl -X POST -H "Content-Type: application/json" \
    -d '{"from": "jed@mail.com", "fromName":"jed","to":"molly@molly.com","subject":"a simple hello", "messageBody":{"errorMessage":"hi!","url":"www.com"}}' \
    http://$API_HOST:$API_PORT/send
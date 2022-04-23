#!/bin/bash

DATA=$(echo "message" |base64)

JSON_STRING=$(jq -n \
                    --arg d $DATA \
                    '{"message":{"attributes":{"teamName":"coucou","projectId":"myProject","eventType":"NEW_QUERIES"},"data":$d,"messageId": "123456"},"subscription":"subscription/path"}')

echo $JSON_STRING

echo "${JSON_STRING}" | http POST localhost:8080/events
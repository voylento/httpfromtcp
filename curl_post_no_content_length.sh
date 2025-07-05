#!/bin/bash

# Send HTTP POST request with JSON body to localhost server
curl -X POST http://localhost:42069/of/httpfromtcp \
  -H "Accept: */*" \
  -H "Content-Type: application/json" \
  -H "User-Name: Bruce" \
  -H "User-Name: Steve" \
  -H "User-Name: Nils" \
  -d '{
    "message": "Hello from curl",
    "timestamp": "2025-07-04",
    "test": true
  }'

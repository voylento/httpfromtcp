#!/bin/bash

# Send HTTP request to localhost server
curl -X GET http://localhost:42069/empty/headers \
  -H "X-Header:" \
  -H "User-Name:"

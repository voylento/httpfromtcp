#!/bin/bash

# Send HTTP request to localhost server
curl -X GET http://localhost:42069/this/is/a/test/of/httpfromtcp \
  -H "User-Agent: curl/7.81.0" \
  -H "Accept: */*" \
  -H "X-Header: Testity-test-test" \
  -H "User-Name: Bruce" \
  -H "User-Name: Steve" \
  -H "User-Name: Nils"

#!/bin/bash
# ENV Variables expected:
# MESHERY_CLOUD_API_URL: Base URL of the Layer5 Cloud API
# MESHERY_TOKEN: Meshery Cloud API authentication token
# EMAIL: Recipient email address
# IMAGE_URI: The URI of the snapshot image


response=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$MESHERY_CLOUD_API_URL/api/integrations/snapshot/email" \
  -H 'Accept: application/json' \
  -H 'Connection: close' \
  -H 'Content-Type: application/json' \
  -H "Cookie: provider_token=$MESHERY_TOKEN" \
  --data "{ \"subject\" : \"Kanvas Design Snapshot\", \"to\":\"$EMAIL\", \"image_uri\" :\"$IMAGE_URI\" }")


# Check HTTP status code
if [ "$response" -eq 200 ]; then
  echo "Email sent"
else
  echo "Failed to send email. HTTP response code: $response"
fi

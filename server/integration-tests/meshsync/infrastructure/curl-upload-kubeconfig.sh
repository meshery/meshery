#!/usr/bin/env bash

# as it is mounted in volume
KUBE_CONFIG_FILE_PATH="/configs/kubeconfig.yaml"

curl -X POST 'http://meshery:9081/api/system/kubernetes' \
  -H 'Accept: */*' \
  -H 'Accept-Language: en-US,en;q=0.9,ru;q=0.8,uk;q=0.7,de;q=0.6' \
  -H 'Connection: keep-alive' \
  -H 'Cookie: meshery-provider=None;' \
  -H 'Origin: http://localhost:9081' \
  -H 'Sec-Fetch-Dest: empty' \
  -H 'Sec-Fetch-Mode: cors' \
  -H 'Sec-Fetch-Site: same-origin' \
  -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36' \
  -H 'sec-ch-ua: "Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "Linux"' \
  -F 'k8sfile=@'"$KUBE_CONFIG_FILE_PATH"';type=application/octet-stream'
#!/bin/sh

cat <<EOF > counting-consul.rendered.json
{
  "Name": "counting",
  "Tags": [
    "v0.0.4"
  ],
  "Address": "${POD_IP}",
  "Port": 9001,
  "Check": {
    "Method": "GET",
    "HTTP": "http://${POD_IP}:9001/health",
    "Interval": "1s"
  }
}
EOF

curl \
    --request PUT \
    --data @counting-consul.rendered.json \
    "http://$HOST_IP:8500/v1/agent/service/register"

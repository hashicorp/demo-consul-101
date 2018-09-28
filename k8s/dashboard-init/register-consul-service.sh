#!/bin/sh

cat <<EOF > dashboard-consul.rendered.json
{
  "Name": "dashboard",
  "Tags": [
    "v0.0.2"
  ],
  "Port": 9002,
  "Check": {
    "Method": "GET",
    "HTTP": "http://${POD_IP}:9002/health",
    "Interval": "1s"
  }
}
EOF

curl \
    --request PUT \
    --data @dashboard-consul.rendered.json \
    "http://$HOST_IP:8500/v1/agent/service/register"

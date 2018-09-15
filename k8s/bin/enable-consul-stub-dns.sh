#!/usr/bin/env sh

# Loads a ConfigMap that allows pods to use the `.consul` TLD.
# Call with the name of your DNS service as deployed by the Consul Helm chart.
#
#     bin/enable-consul-stub-dns.sh piquant-shark-consul-dns
#
# https://www.consul.io/docs/platform/k8s/dns.html

DNS_SERVICE_NAME=$1

cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  labels:
    addonmanager.kubernetes.io/mode: EnsureExists
  name: kube-dns
  namespace: kube-system
data:
  stubDomains: |
    {"consul": ["$(kubectl get svc $DNS_SERVICE_NAME -o jsonpath='{.spec.clusterIP}')"]}
EOF

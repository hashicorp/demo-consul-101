#!/usr/bin/env sh
# Copyright IBM Corp. 2018, 2023
# SPDX-License-Identifier: MPL-2.0


# Loads a ConfigMap that allows pods to use the `.consul` TLD.
# https://www.consul.io/docs/platform/k8s/dns.html

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
    {"consul": ["$(kubectl get svc consul-dns -o jsonpath='{.spec.clusterIP}')"]}
EOF

# Kubernetes Services with Consul Service Discovery

IN PROCESS

## Prerequisites

You must create a GCP IAM service account and authenticate with it on the command line.

We also assume that you have installed the `gcloud` command line tool and `helm`.

### Service Account Details

https://cloud.google.com/sdk/gcloud/reference/auth/activate-service-account

    gcloud auth activate-service-account --key-file="my-consul-service-account.json"

## Configure kubectl for cluster

Go to the web UI, find cluster. Click `Connect` button. Copy the snippet and paste it into your terminal.

    gcloud container clusters get-credentials my-consul-cluster \
      --zone us-west1-b --project my-project

## Helm

Install Helm to the k8s cluster.

```sh
$ helm init

Tiller (the Helm server-side component) has been installed into your Kubernetes Cluster.

Please note: by default, Tiller is deployed with an insecure 'allow unauthenticated users' policy.
To prevent this, run `helm init` with the --tiller-tls-verify flag.
For more information on securing your installation see: https://docs.helm.sh/using_helm/#securing-your-helm-installation
Happy Helming!
```

Create permissions for the service account.

```sh
kubectl create clusterrolebinding add-on-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:default
```

Edit values in `helm-consul-values.yml` if desired.

```yaml
global:
  datacenter: hashidc1

ui:
  service:
    type: "LoadBalancer"
```

Install Consul to the cluster, either from the stable repository or from the development [GitHub repo](https://github.com/hashicorp/consul-helm).

```sh
helm install --name consul-release -f helm-consul-values.yaml stable/consul
```

## Enable stub-dns

https://www.consul.io/docs/platform/k8s/dns.html

Find the name of your `dns` service with

```sh
$ kubectl get svc
```

Pass the service name to the stub dns script.

```sh
$ bin/enable-consul-stub-dns.sh lucky-penguin-consul-dns
```

## Apply the resources

Deploy with Kubernetes defined by all files in the `yaml` directory.

```sh
$ kubectl apply -f yaml/
```

Refresh your [GCP](https://console.cloud.google.com/kubernetes) console. Go to "Services" and you should see a public IP address for the `dashboard-service-deployment-service`. Visit it to see the dashboard and counting service which are communicating to each other using Consul service discovery.

## Debugging

    kubectl exec -it my-pod-name /bin/sh

    apk add curl

    kubectl logs my-pod-name

## Advanced

Scale up deployments to start more counting services.

    kubectl get deployments
    kubectl scale deployments/counting-service-deployment --replicas=2

## Other/Random Notes

### Run Service

https://kubernetes.io/docs/concepts/overview/object-management-kubectl/declarative-config/

https://docs.docker.com/get-started/part2/

https://www.consul.io/api/agent/service.html#register-service

### Start K8S Dashboard

    kubectl proxy --port=8080

### Get Bearer Token for K8S Dashboard

    gcloud config config-helper --format=json

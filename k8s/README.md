# Kubernetes Services with Consul Service Discovery

IN PROCESS

## Prerequisites and Setup

We assume that you have installed the `gcloud` command line tool, `helm`, and `kubectl`.

https://cloud.google.com/sdk/docs/downloads-interactive

```sh
$ gcloud init
```

Install `helm` and `kubectl` with Homebrew.

```sh
$ brew install kubernetes-cli
$ brew install kubernetes-helm
```

### Service account authentication

It's recommended that you create a GCP IAM service account and authenticate with it on the command line.

https://console.cloud.google.com/iam-admin/serviceaccounts

https://cloud.google.com/sdk/gcloud/reference/auth/activate-service-account

```sh
$ gcloud auth activate-service-account --key-file="my-consul-service-account.json"
```

### Create a kubernetes cluster

https://console.cloud.google.com/kubernetes/list

Click "Create Cluster" and use the defaults. Find the "Create" button at the bottom of the drawer.

### Configure kubectl to talk to your cluster

Go to the web UI, find "Clusters". Click the "Connect" button. Copy the snippet and paste it into your terminal.

```sh
$ gcloud container clusters get-credentials my-consul-cluster \
      --zone us-west1-b --project my-project
```

## Install helm to your cluster

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

Create a new file named `helm-consul-values.yml`. Edit to expose a load balancer so you can view the Consul UI.

```yaml
global:
  datacenter: hashidc1

ui:
  service:
    type: "LoadBalancer"

syncCatalog:
  enabled: true
```

Install Consul to the cluster, either from the stable repository or from the development [GitHub repo](https://github.com/hashicorp/consul-helm).

```sh
helm install -f helm-consul-values.yaml stable/consul
```

## Enable stub-dns

https://www.consul.io/docs/platform/k8s/dns.html

Find the name of your `dns` service with

```sh
$ kubectl get svc
```

Pass the service name to the stub DNS script in this demo repo.

```sh
$ bin/enable-consul-stub-dns.sh lucky-penguin-consul-dns
```

## Apply the resources

Deploy an application with Kubernetes. Use all files in the `yaml` directory.

```sh
$ kubectl apply -f yaml/
```

Refresh your [GCP](https://console.cloud.google.com/kubernetes) console. Go to "Services" and you should see a public IP address for the `dashboard-service-load-balancer`. Visit it to see the dashboard and counting service which are communicating to each other using Consul service discovery.

## Extra: Debugging

```sh
# Connect to a container
kubectl exec -it my-pod-name /bin/sh

# Install tools on a container for curl and dig
apk add curl
apk add bind-tools

# View logs for a pod
kubectl logs my-pod-name

# See full configuration for debugging
helm template stable/consul
```

## Advanced

Scale up deployments to start more counting services.

```sh
kubectl get deployments
kubectl scale deployments/counting-service-deployment --replicas=2
```

## Other/Random Notes

### Run Service

https://kubernetes.io/docs/concepts/overview/object-management-kubectl/declarative-config/

https://docs.docker.com/get-started/part2/

https://www.consul.io/api/agent/service.html#register-service

### Start K8S Dashboard

    kubectl proxy --port=8080

### Get Bearer Token for K8S Dashboard

    gcloud config config-helper --format=json

### Catalog sync

These permissions may be needed:

```sh
$ kubectl set subject clusterrolebinding system:node --group=system:nodes

kubectl create rolebinding admin --clusterrole=admin --user="system:serviceaccount:default:default" --namespace=default
```

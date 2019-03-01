# Kubernetes Services with Consul Service Discovery

This README briefly sketches the sequence of steps for setting up Google Cloud for Kubernetes and Consul.

For a Microsoft Azure version, see the README and Terraform configs in the [terraform/azure](https://github.com/hashicorp/demo-consul-101/tree/master/k8s/terraform/azure) directory.

For either, you can use the same YAML files to deploy the demo applications to the Kubernetes cluster.

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

## Task 1: Run the simplest container

In this task, you'll deploy a container that returns JSON that includes a number and the name of the host. You'll put a load balancer in front so you can see the output.

### Create minimal yaml config

Here is the simplest configuration that will deploy a container. Create a file named `counting-minimal.yaml` and paste these contents into it.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: counting-minimal-pod
  labels:
    app: counting-service
spec:
  containers:
    - name: counting-service
      image: topfunky/counting-service:0.0.1
      ports:
      - containerPort: 9001
```

This code has been created in this repository. Apply to the cluster with:

```sh
$ kubectl apply -f 01-yaml-minimal

pod/counting-minimal-pod created
```

Look at the Google web console. look at logs. You should see

```plaintext
Serving at http://localhost:9001
```

Connect to the pod from your local machine.

```sh
$ kubectl port-forward pod/counting-minimal-pod 9001:9001
```

Now visit http://localhost:9001/

You should see JSON that contains a number and the name of the host.

### Implement load balancer

Sometimes we want a pod to surface on an IP address outside of the cluster. Let's add a load balancer.

Add the following under the pod definition in the same `counting-minimal.yaml` file (or see the completed code in this repository).

```yaml
---
apiVersion: v1
kind: Service
metadata:
  name: counting-minimal-load-balancer
spec:
  ports:
  - protocol: "TCP"
    port: 80
    targetPort: 9001
  selector:
    app: counting-service
  type: LoadBalancer

```

Apply with

```sh
$ kubectl apply -f yaml-minimal/counting-minimal.yaml

pod/counting-minimal-pod unchanged
service/counting-minimal-load-balancer created
```

Visit the Google Cloud console and go to "Services." You will see a service of type "Load Balancer" with an IP address next to it. Click the IP address and you'll see JSON output from the counting service.

```json
{"count":1,"hostname":"counting-minimal-pod"}
```

### Gather data

```sh
$ kubectl get pods
```

```sh
$ kubectl logs counting-minimal-pod

Serving at http://localhost:9001
(Pass as PORT environment variable)
```

```sh
$ kubectl get pods --output=json
```

Connect to a running pod.

```sh
$ kubectl exec -it counting-minimal-pod /bin/sh
```

Run commands like `env` or try to start `counting-service` manually with a different port.

```sh
$ PORT=9002 ./counting-service
```

Because this container is built from Alpine Linux, it doesn't have many development tools. You can install them:

```sh
$ apk add curl
$ apk add bind-tools
```

### Delete the pod

We're done with this pod, so delete it (and the load balancer).

```sh
$ kubectl delete -f yaml-minimal/counting-minimal.yaml
```

## Task 2: Install helm to your cluster

Install Helm to the k8s cluster.

```sh
$ helm init

Tiller (the Helm server-side component) has been installed into your Kubernetes Cluster.

Please note: by default, Tiller is deployed with an insecure 'allow unauthenticated users' policy.
To prevent this, run `helm init` with the --tiller-tls-verify flag.
For more information on securing your installation see: https://docs.helm.sh/using_helm/#securing-your-helm-installation
Happy Helming!
```

Go to the Google Cloud console and choose "Services." Select "Show System Objects." You should see an object named `tiller-deploy`. This is the server side component of Helm.

Next, create permissions for the service account so it can install Helm charts (packages).

```sh
$ kubectl create clusterrolebinding add-on-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:default

clusterrolebinding.rbac.authorization.k8s.io "add-on-cluster-admin" created
```

Create a new file named `helm-consul-values.yaml`. Edit to expose a load balancer so you can view the Consul UI across the internet.

```yaml
global:
  datacenter: hashidc1

ui:
  service:
    type: "LoadBalancer"

syncCatalog:
  enabled: true
```

Install Consul to the cluster. We'll use a clone of the development [GitHub repo](https://github.com/hashicorp/consul-helm).

```sh
$ git clone https://github.com/hashicorp/consul-helm.git

$ helm install -f helm-consul-values.yaml ./consul-helm
```

Verify that this worked by going to "Services" in the Google Cloud console. Find the load balancer for `*-consul-ui`. Click the IP address and you'll see the Consul web UI.

## Task 3: Enable stub-dns

In order for Consul service discovery to work smoothly, we need to enable Consul within the Kubernetes DNS system.

https://www.consul.io/docs/platform/k8s/dns.html

Find the name of your `dns` service with

```sh
$ kubectl get svc
```

Pass the service name matching `*-consul-dns` to the stub DNS script in this demo repo.

```sh
$ bin/enable-consul-stub-dns.sh lucky-penguin-consul-dns

configmap "kube-dns" configured
```

## Task 4: Use Consul K/V

Let's use Consul's key/value store.

Get a list of pods and find one that is running a Consul agent. We'll use this as an easy way to run Consul CLI commands.

```sh
$ kubectl get pods
```

Look for one with `consul` in the name. Connect to the running pod.

```sh
$ kubectl exec -it giggly-echidna-consul-5t2dc /bin/sh
```

Once connected, run a command that saves a value to Consul.

```sh
$ consul kv put redis/config/connections 5
```

Go to the Consul web UI and look in the **Key/Value** tab. You should see the hierarchy that contains the `redis/config/connections` value.

Use the web UI to change the value to 10. Go back to the pod and `get` the value.

```sh
$ consul kv get redis/config/connections

10
```

You're now running Consul's key/value store and can work with data.

-> **NOTE:** Neither `envconsul` or `consul-template` are installed in this container and must be installed separately if you plan to use them.

## Task 5: Create an application that uses Consul service discovery

Deploy an application with Kubernetes. The yaml files for deploying this application are in this repository. 

```sh
$ kubectl create -f 02-yaml-discovery/
```

Refresh your [GCP](https://console.cloud.google.com/kubernetes) console. Go to "Services" and you should see a public IP address for the `dashboard-service-load-balancer`. Visit it to see the dashboard and counting service which are communicating to each other using Consul service discovery. (See code in `dashboard-service` for details.)

## Task 6: Create an application that uses Consul Connect secure service segmentation

The `counting` service needs to start an extra container running `consul` that manually starts its own proxy. The consul binary can be found on your consul pod. The following pod is the same example from above.

```sh
kubectl exec -it giggly-echidna-consul-5t2dc /bin/sh
```

```sh
$ exec /bin/consul connect proxy \
      -http-addr=${HOST_IP}:8500 \
      -service=counting \
      -service-addr=127.0.0.1:9001 \
      -listen=${POD_IP}:19001 \
      -register
```

The `dashboard` service needs to start an extra container running `consul` that manually starts an upstream proxy to the `counting` service proxy.

```sh
$ exec /bin/consul connect proxy \
  -http-addr=${HOST_IP}:8500 \
  -service=dashboard \
  -upstream="counting:9001"
```

## Extra: Debugging

```sh
# Connect to a container
$ kubectl exec -it my-pod-name /bin/sh

# View logs for a pod
$ kubectl logs my-pod-name

# See full configuration for debugging
$ helm template stable/consul
```

Within a pod (may require Consul pod or extra installation of `curl`).

```sh
# View all environment variables
$ env

# Install tools on a container for curl and dig
$ apk add curl
$ apk add bind-tools

# Use the Consul HTTP API from any pod
$ curl http://consul.service.consul:8500/v1/catalog/datacenters

# Use service discovery
$ ping dashboard.service.consul
```

## Advanced

Scale up deployments to start more counting services.

```sh
$ kubectl get deployments
$ kubectl scale deployments/counting-service-deployment --replicas=2
```

Or in a deployment:

```yaml
spec:
  replicas: 5
```

Health checks:

```yaml
spec:
  containers:
    - name: "..."
      livenessProbe:
        # an http probe
        httpGet:
          path: /health
          port: 9002
        # length of time to wait for a pod to initialize
        # after pod startup, before applying health checking
        initialDelaySeconds: 30
        timeoutSeconds: 1
      # Other content omitted
```

https://kubernetes.io/docs/tutorials/k8s101/

https://kubernetes.io/docs/tutorials/k8s201/

Networking: https://kubernetes.io/docs/tutorials/services/source-ip/

```sh
$ kubectl get nodes --out=yaml
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

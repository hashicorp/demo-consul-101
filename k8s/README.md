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

In this task, you'll deploy a container that returns JSON that includes a number and the name of the host. You'll put a
load balancer in front so that you can see the output.

### Create minimal yaml config

Here is the simplest configuration that will deploy a container. Create a file named `counting-minimal.yaml` and paste these
contents into it.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: counting-minimal-pod
  labels:
    app: counting
spec:
  containers:
    - name: counting
      image: hashicorp/counting-service:0.0.1
      ports:
      - containerPort: 9001
```

This code has been created in this repository. Apply to the cluster with:

```sh
$ kubectl apply -f 01-yaml-minimal/counting-service.yaml

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
    app: counting
  type: LoadBalancer

```

Apply with

```sh
$ kubectl apply -f 01-yaml-minimal/counting-load-balancer.yaml

service/counting-minimal-load-balancer created
```

Visit the Google Cloud console and go to "Services." You will see a service of type "Load Balancer" with an IP address
next to it. Click the IP address and you'll see JSON output from the counting service.

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

### Delete the pod

We're done with this pod, so delete it (and the load balancer).

```sh
$ kubectl delete -f 01-yaml-minimal
```

## Task 2: Use helm to install Consul to your cluster

Now you will download the official Consul Helm chart and use it to install Consul. To install the chart, you will need to
add the hashicorp Helm repo using `helm repo add`.

```sh
$ helm repo add hashicorp https://helm.releases.hashicorp.com
```

Review the file in this directory named `config.yaml`. It has a minimal Consul configuration that will allow
you to get started using Consul.

```yaml
global:
  datacenter: dc1
  
ui:
  service:
    type: "LoadBalancer"

controller:
  enabled: true

connectInject:
  enabled: true
```

Install Consul to the cluster using he values file.

```
$ helm install -f config.yaml ./consul-helm --version "0.30.0"
```

Verify that this worked by going to "Services" in the Google Cloud console. Find the load balancer for `*-consul-ui`.
Click the IP address and you'll see the Consul web UI.

## Task 3: Use Consul K/V

Let's use Consul's key/value store.

Connect to the `consul-server-0` pod.

```sh
$ kubectl exec -it consul-server-0 -- /bin/sh
```

Once connected, run a command that saves a value to Consul.

```sh
$ consul kv put redis/config/connections 5
Success! Data written to: redis/config/connections
```

Go to the Consul web UI and look in the **Key/Value** tab. You should see the hierarchy that contains the `redis/config/connections` value.

Use the web UI to change the value to 10. Go back to the pod and `get` the value.

```sh
$ consul kv get redis/config/connections
10
```

You're now running Consul's key/value store and can work with data.

Exit from the `consul-server-0` pod to continue to the next task.

```
$ exit
```

-> **NOTE:** Neither `envconsul` or `consul-template` are installed in this container and must be installed separately
if you plan to use them.

## Task 4: Create an application that uses Consul service discovery

In order for Consul service discovery to work, we need to enable Consul within the Kubernetes DNS system. See the
[documentation](https://www.consul.io/docs/platform/k8s/dns.html) for details.

Run the stub DNS script in this demo repo. This script will add a [Stub-domain](https://kubernetes.io/docs/tasks/administer-cluster/dns-custom-nameservers/)
config map entry for the Consul DNS server.

```sh
$ bin/enable-consul-stub-dns.sh
configmap/kube-dns configured
```

Deploy an application with Kubernetes. The yaml files for deploying this application are in this repository. Notice that
the `dashboard-server` uses `counting.service.consul` to reach the counting service. This is Consul DNS enabling service
discovery!

```sh
$ kubectl apply -f 02-yaml-discovery/
pod/counting-deployment created
pod/dashboard created
service/dashboard-load-balancer created
```

Refresh your [GCP](https://console.cloud.google.com/kubernetes) console. Go to "Services", and you should see a public IP
address for the `dashboard-service-load-balancer`. Visit it to see the dashboard and counting service which are communicating
to each other using Consul service discovery. (See code in `dashboard-service` for details.)

## Task 5: Create an application that uses Consul service mesh secure service segmentation

The `counting` service needs to start an extra container running `consul` that manually starts its own proxy. The consul
binary can be found on your consul pod. The following pod is the same example from above.

```sh
kubectl apply -f 04-yaml-connect-envoy
```

Retrieve the public IP address of the `dashboard-load-balancer`. If you the `EXTERNAL-IP` for the load balancer is set to
`<pending>` wait a minute or two and try again. You may have to run this several times until the resource is allocated
by your cloud provider.

```sh
$ kubectl get svc
NAME                              TYPE           CLUSTER-IP    EXTERNAL-IP     PORT(S)                                                                   AGE
consul-connect-injector-svc       ClusterIP      10.72.1.252   <none>          443/TCP                                                                   7d23h
consul-controller-webhook         ClusterIP      10.72.10.41   <none>          443/TCP                                                                   7d23h
consul-dns                        ClusterIP      10.72.5.39    <none>          53/TCP,53/UDP                                                             7d23h
consul-ingress-gateway            LoadBalancer   10.72.1.11    35.225.10.155   8080:32306/TCP,8443:32079/TCP                                             7d22h
consul-server                     ClusterIP      None          <none>          8500/TCP,8301/TCP,8301/UDP,8302/TCP,8302/UDP,8300/TCP,8600/TCP,8600/UDP   7d23h
consul-ui                         ClusterIP      10.72.1.246   <none>          80/TCP                                                                    7d23h
dashboard-load-balancer           LoadBalancer   10.72.1.67    34.66.3.28      80:30840/TCP                                                              3m19s
kubernetes                        ClusterIP      10.72.0.1     <none>          443/TCP                                                                   8d
```

## Extra: Debugging

Start an interactive terminal on a pod running in the cluster.

```sh
$ kubectl exec -it consul-server-0 /bin/sh
```

View logs for a pod.

```
$ kubectl logs consul-server-0
```

Review the helm chart configuration that was applied during installation.

```
$ helm template hashicorp/consul
```

Inspect environment variabels from within a pod.

```sh
$ env
```

Validate service discovery

```
$ kubectl exec dashboard -- ping counting.service.consul
```

## Advanced

If you are using deployments,  you can scale them up to start more service instances using `kubectl scale`.

```sh
$ kubectl get deployments
$ kubectl scale deployments/counting --replicas=2
```

Or in a deployment definition within a yaml file.

```yaml
spec:
  replicas: 5
```

You can also add health checks to your pod or deployment specs.

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

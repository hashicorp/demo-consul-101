## Azure AKS (Simple)

```
az ad sp create-for-rbac --skip-assignment

// Then copy appId and password to terraform.tfvars

// Provision resource group, AKS cluster
cd 01-create-aks-cluster
terraform init
terraform apply

// Provision tiller service account, install helm
cd ../02-fix-k8s-rbac
terraform init
terraform apply
cd ..

// Verify that tiller is installed
kubectl get pods --all-namespaces | grep tiller

// View k8s dashboard
az aks browse --resource-group geoffrey-rg --name geoffrey-aks

git clone https://github.com/hashicorp/demo-consul-101.git
cd demo-consul-101/k8s
git clone https://github.com/hashicorp/consul-helm.git

// View helm-consul-values.yaml

helm install -f helm-consul-values.yaml --name=azure ./consul-helm
// View Consul UI
kubectl get service azure-consul-ui --watch

// Deploy application
kubectl apply -f 04-yaml-connect-envoy

// Then configure intentions
// Then show container configuration
```

## Azure AKS (Full)

```
az ad sp create-for-rbac --skip-assignment

// Try
// https://docs.microsoft.com/en-us/cli/azure/group?view=azure-cli-latest#az-group-update
// Find AKS-specific roles to add to service principal
az group update ...
az role assignment ...

{
  "appId": "aaaaaaa",
  "displayName": "azure-cli-2019-01-04-16-44-26",
  "name": "http://azure-cli-2019-01-04-16-44-26",
  "password": "aaaaaaaa",
  "tenant": "aaaaaa"
}

APP_ID="..."
APP_PWD="..."
MY_RG="geoffreyRG"
MY_CLUSTER="geoffreyAKSCluster"

// Create resource group geoffreyRG
az group create --location westus2 --name geoffreyRG

az aks create \
    --resource-group $MY_RG \
    --name $MY_CLUSTER \
    --node-count 3 \
    --service-principal $APP_ID \
    --client-secret $APP_PWD \
    --generate-ssh-keys

// Or without service principal (will be created)

az aks create \
    --resource-group $MY_RG \
    --name $MY_CLUSTER \
    --node-count 3 \
    --generate-ssh-keys

az aks get-credentials --resource-group $MY_RG --name $MY_CLUSTER

kubectl create clusterrolebinding kubernetes-dashboard -n kube-system --clusterrole=cluster-admin --serviceaccount=kube-system:kubernetes-dashboard

// Accounts for helm and tiller
kubectl create serviceaccount --namespace kube-system tiller
kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
helm init --service-account tiller

az aks browse --resource-group $MY_RG --name $MY_CLUSTER

git clone https://github.com/hashicorp/demo-consul-101.git
cd demo-consul-101/k8s

helm init
kubectl get pods --all-namespaces | grep tiller
git clone https://github.com/hashicorp/consul-helm.git

helm install -f helm-consul-values.yaml ./consul-helm
kubectl get service consul-ui --watch

kubectl apply -f 04-yaml-connect-envoy
```
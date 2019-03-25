# Demo Consul 101

Demo code and microservices for the HashiCorp Consul 101 course.

Email training@hashicorp.com or see https://www.hashicorp.com/training for details.

## Quickstart: Docker Compose

To run both microservices with Docker Compose (but without Consul), run `docker compose up`.

```
$ cd demo-consul-101
$ docker-compose up
```

You can view the operational application dashboard at http://localhost:8080

A subsequent evolution of the application would involve registering each service with Consul and using Consul DNS to configure services to discover each other.

## Quickstart: Consul Connect

More documentation is coming. In the meantime, you can start a local demo with:

```
consul agent -dev -config-dir="./demo-config-localhost" -node=laptop
```

Then start instances of `dashboard-service` and `counting-service`

```
cd services/dashboard-service
PORT=9002 go run main.go

cd services/counting-service
PORT=9003 go run main.go

cd services/counting-service
PORT=9004 go run main.go

consul connect proxy -sidecar-for counting-1
consul connect proxy -sidecar-for counting-2

consul connect proxy -sidecar-for dashboard
```

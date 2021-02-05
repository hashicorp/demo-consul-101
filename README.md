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

A subsequent evolution of the application would involve registering each service with Consul and using Consul DNS to
configure services to discover each other.

## Quickstart: Consul service mesh

If you have the Consul binary installed locally, you can use the following sequence of commands to run a demo mesh on
your local laptop.

Start a local Consul dev agent with:

```
consul agent -dev -config-dir="./demo-config-localhost" -node=laptop
```

Start the `dashboard-service` in a separate shell session.

```
PORT=9002 go run ./services/dashboard-service/main.go
```

Start the `counting-service` in a separate shell session.

```
PORT=9003 go run ./services/counting-service/main.go
```

Start a second instance of the `counting-service` in a separate shell session.

```
PORT=9004 go run ./services/counting-service/main.go
```

Start the sidecar proxy for the `counting-1` service instance in a separate shell session.

```
consul connect proxy -sidecar-for counting-1
```

Start the sidecar proxy for the `counting-2` service instance in a separate shell session.

```
consul connect proxy -sidecar-for counting-2
```

Start the sidecar proxy for the `dashboard` service in a separate shell session.

```
consul connect proxy -sidecar-for dashboard
```

Now visit the application at `localhost:9002`.

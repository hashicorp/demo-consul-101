# Dashboard Service

A web application that displays a numeric dashboard. It retrieves a count from a backend counting service and displays a live update.

For use in learning to use Consul for service discovery and segmentation (connection via secure proxies).

Defaults to running on port 80. Set `PORT` as ENV var to specify another port.

Defaults to looking for the `counting-service` running at `localhost:9001`. Can be set with the `COUNTING_SERVICE_URL` ENV var.

### Run precompiled binary

To run with the defaults (port 80, looking for the backend counting service at `localhost:9001`):

    dashboard-service

To run on a specific port or looking for the `counting-service` elsewhere:

    PORT=9002 COUNTING_SERVICE_URL=counting.service.consul dashboard-service

### Build

Build for Linux and Darwin:

    ./bin/build

Output can be found in `dist`.

### Run from source

    go get
    PORT=9002 go run main.go

### View

    http://localhost:9002

### Dependencies

This application assumes that a counting service is running on `localhost:9001`.

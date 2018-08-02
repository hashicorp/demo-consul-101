# Dashboard Service

A web application that displays a numeric dashboard. It retrieves a count from a backend counting service and displays a live update.

For use in learning to use Consul for service discovery and segmentation (connection via secure proxies).

Defaults to running on port 80. Set `PORT` as ENV var to specify another port.

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


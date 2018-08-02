# Counting Service

A web service that increments a count every time it is invoked. It returns the current count as JSON.

For use in learning to use Consul for service discovery and segmentation (connection via secure proxies).

Defaults to running on port 80. Set `PORT` as ENV var to specify another port.

### Build

Build for Linux and Darwin:

    ./bin/build

Output can be found in `dist`.

### Run from source

    go get
    PORT=9001 go run main.go

### View

    http://localhost:9001

### Dependencies

This application is intended to be used by the dashboard service.

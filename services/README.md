# Demo Consul Services

These services are simple workloads designed to use in labs and demonstrations
of Consul, Consul Connect, and Nomad.  We build them for a large number of
target platforms to enable demonstrations over heterogeneous platforms.

## Building the service packages

Requires an installed and configured golang development environment

- **Bootstrap your build environment with `./bootstrap`.** - This will install
  the necessary go dependencies 

- **Run `./build_all`.** - This will create all of the executables in a `pkg`
  directory. It will zip and calculate a SHA256 sum of the zipfiles into a `zip`
  directory. The default configuration will build a tree that looks like:

    ```
    pkg
    |-- counting-service_darwin_amd64
    |-- counting-service_freebsd_386
    |-- counting-service_freebsd_amd64
    |-- counting-service_freebsd_arm
    |-- counting-service_linux_386
    |-- counting-service_linux_amd64
    |-- counting-service_linux_arm
    |-- counting-service_linux_arm64
    |-- counting-service_solaris_amd64
    |-- counting-service_windows_386.exe
    |-- counting-service_windows_amd64.exe
    |-- dashboard-service_darwin_amd64
    |-- dashboard-service_freebsd_386
    |-- dashboard-service_freebsd_amd64
    |-- dashboard-service_freebsd_arm
    |-- dashboard-service_linux_386
    |-- dashboard-service_linux_amd64
    |-- dashboard-service_linux_arm
    |-- dashboard-service_linux_arm64
    |-- dashboard-service_solaris_amd64
    |-- dashboard-service_windows_386.exe
    |-- dashboard-service_windows_amd64.exe
    +-- zips
        |-- SHA256SUMS.txt
        |-- counting-service_darwin_amd64.zip
        |-- counting-service_freebsd_386.zip
        |-- counting-service_freebsd_amd64.zip
        |-- counting-service_freebsd_arm.zip
        |-- counting-service_linux_386.zip
        |-- counting-service_linux_amd64.zip
        |-- counting-service_linux_arm.zip
        |-- counting-service_linux_arm64.zip
        |-- counting-service_solaris_amd64.zip
        |-- counting-service_windows_386.exe.zip
        |-- counting-service_windows_amd64.exe.zip
        |-- dashboard-service_darwin_amd64.zip
        |-- dashboard-service_freebsd_386.zip
        |-- dashboard-service_freebsd_amd64.zip
        |-- dashboard-service_freebsd_arm.zip
        |-- dashboard-service_linux_386.zip
        |-- dashboard-service_linux_amd64.zip
        |-- dashboard-service_linux_arm.zip
        |-- dashboard-service_linux_arm64.zip
        |-- dashboard-service_solaris_amd64.zip
        |-- dashboard-service_windows_386.exe.zip
        +-- dashboard-service_windows_amd64.exe.zip
    ```

- **Clean up build artifacts with `./clean_all`

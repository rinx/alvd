<div align="center">
<img src="./assets/cover.svg" width="100%">
</div>

alvd - A Lightweight Vald
===

[![License: Apache 2.0](https://img.shields.io/github/license/rinx/alvd.svg?style=flat-square)](https://opensource.org/licenses/Apache-2.0)
[![release](https://img.shields.io/github/release/rinx/alvd.svg?style=flat-square)](https://github.com/rinx/alvd/releases/latest)
[![ghcr.io](https://img.shields.io/badge/ghcr.io-rinx%2Falvd-brightgreen?logo=docker&style=flat-square)](https://github.com/users/rinx/packages/container/package/alvd)
[![Docker Pulls](https://img.shields.io/docker/pulls/rinx/alvd.svg?style=flat-square)](https://hub.docker.com/r/rinx/alvd)

A lightweight distributed vector search engine based on [Vald](https://vald.vdaas.org) codebase.

- works without Kubernetes
- single binary (less than 30MB)
- easy to run (can be configured by command-line options)
- consists of Agent and Server
    - alvd has almost same features that Vald's gateway-lb + discoverer and agent-ngt have.

alvd is highly inspired by [k3s](https://k3s.io) project.

Rationale
---

Vald is an awesome highly scalable distributed vector search engine works on Kubernetes.
It has great features such as file-based backup, metrics-based ordering of Agents. Also Vald is highly configurable using YAML files.

However it requires

- Kubernetes APIs to discover Vald Agents
- knowledge of operating Kubernetes
- knowledge of tuning a lot of complicated parameters

it is a little difficult for the users.

In this project, we eliminated several features of Vald such as (meta, backup manager, index manager, etc...) and just focused on Vald's gateway-lb and agent-ngt.
By using [rancher/remotedialer](https://github.com/rancher/remotedialer), Vald's discoverer feature is not needed anymore.
Also we eliminated advanced options and adopt command-line options for configuring the application behavior instead of YAML files.

As stated above, alvd is focused on "easy to use", "Kubernetes-less" and "less components".

How it works
---

<div align="center">
<a href="./assets/howitworks.svg">
<img src="./assets/howitworks.svg" width="100%">
</a>
</div>

Quick Start
---

1. Get the latest build from [Release](https://github.com/rinx/alvd/releases) page and unzip it.
2. Run alvd Server.
    ```sh
    $ ./alvd server
    2020-12-18 19:18:27     [INFO]: start alvd server
    2020-12-18 19:18:27     [INFO]: metrics server starting on 0.0.0.0:9090
    2020-12-18 19:18:27     [INFO]: websocket server starting on 0.0.0.0:8000
    2020-12-18 19:18:27     [INFO]: gateway gRPC API starting on 0.0.0.0:8080
    INFO[0000] Connecting to proxy                           url="ws://0.0.0.0:8000/connect"
    2020-12-18 19:18:27     [INFO]: agent gRPC API starting on 0.0.0.0:8081
    INFO[0000] Handling backend connection request [7q6ai4gbve83spij0s4g]
    2020-12-18 19:18:27     [INFO]: connected to: 0.0.0.0:8000
    ```
    alvd Server's websocket server starts on 0.0.0.0:8000 and alvd Server's gRPC API starts on 0.0.0.0:8080.
    Also, alvd Agent's gRPC API starts on 0.0.0.0:8081 (alvd Agent process on the Server can be disabled using `--agent=false` option).
3. Run alvd Agent on a different node (or a different terminal on the same node with `--server 0.0.0.0:8000` and `--grpc-port 8082` option).
    ```sh
    $ ./alvd agent --server host-of-server-node:8000
    $ # ./alvd agent --server 0.0.0.0:8000 --grpc-port 8082 --metrics-port=9091
    2020-12-18 19:20:15     [INFO]: start alvd agent
    2020-12-18 19:20:15     [INFO]: metrics server starting on 0.0.0.0:9090
    INFO[0000] Connecting to proxy                           url="ws://host-of-server-node:8000/connect"
    2020-12-18 19:20:15     [INFO]: agent gRPC API starting on 0.0.0.0:8081
    2020-12-18 19:20:15     [INFO]: connected to: host-of-server-node:8000
    ```
4. Add more alvd Agents on the other nodes (or the other ports on the same node).
    ```sh
    $ ./alvd agent --server host-of-server-node:8000
    $ # ./alvd agent --server 0.0.0.0:8000 --grpc-port 808{3,4,5} --metrics-port=909{2,3,4}
    ```
5. Now we can access the alvd Server's gRPC API (`host-of-server-node:8080`) using Vald v1 clients.
    If you don't have one, you can use [valdcli](https://github.com/vdaas/vald-client-clj) (this CLI is built for linux-amd64 and macos-amd64).
    ```sh
    $ # insert 100 vectors (dimension: 784) with random IDs
    $ valdcli rand-vecs -d 784 -n 100 --with-ids | valdcli -h host-of-server-node -p 8080 stream-insert
    $ # search a random vector
    $ valdcli rand-vec -d 784 | valdcli -h host-of-server-node -p 8080 search
    ```

Distribution
---

On [Release](https://github.com/rinx/alvd/releases) page, alvd binaries for amd64 Linux machines are available.

- alvd-linux-amd64.zip doesn't use AVX instructions.
- alvd-linux-amd64-avx2.zip uses AVX2 instuctions for distance calculations. It is faster.

Docker images are available on GitHub Package Registries and DockerHub.
The images tagged by `noavx` are built for amd64, arm64 and armv7 architectures.
`avx2` images are only available for amd64 architectures.

- [ghcr.io/rinx/alvd](https://github.com/users/rinx/packages/container/package/alvd)
- [rinx/alvd](https://hub.docker.com/r/rinx/alvd)

Running on Docker using Docker Compose
---

There's an example [`docker-compose.yml`](https://github.com/rinx/alvd/tree/main/docker-compose.yml) in this repository.

Try to run it with a command:

```sh
$ docker-compose up
Starting alvd-agent-2 ... done
Starting alvd-agent-1 ... done
Starting alvd-agent-3 ... done
Starting alvd-server  ... done
Attaching to alvd-agent-2, alvd-agent-3, alvd-agent-1, alvd-server
alvd-agent-1    | 2021-04-30 02:33:36   [INFO]: start alvd agent
alvd-agent-1    | 2021-04-30 02:33:36   [INFO]: metrics server starting on 0.0.0.0:9090
alvd-agent-2    | 2021-04-30 02:33:35   [INFO]: start alvd agent
alvd-agent-2    | 2021-04-30 02:33:35   [INFO]: metrics server starting on 0.0.0.0:9090
alvd-agent-2    | 2021-04-30 02:33:35   [INFO]: agent gRPC API starting on 0.0.0.0:8081
...
alvd-server     | 2021-04-30 02:33:36   [INFO]: gateway gRPC API starting on 0.0.0.0:8080
...
```

Then 1 Server + 3 Agents will run on your Docker environment.
We can access to the alvd Server's gRPC API on 8080 using Vald v1 clients.

```sh
$ # insert 100 vectors (dimension: 784) with random IDs
$ valdcli rand-vecs -d 784 -n 100 --with-ids | valdcli -h localhost -p 8080 stream-insert
$ # search a random vector
$ valdcli rand-vec -d 784 | valdcli -h localhost -p 8080 search
```

The metrics APIs are exported on 9090-9093 ports. We can access them using curl.

```sh
$ curl http://localhost:9090/metrics
$ curl http://localhost:9091/metrics
...
```

In the `docker-compose.yml` file, there're definitions of Prometheus and Grafana services. If they are enabled, a metrics dashboard can be displayed on your machine. (http://localhost:3000)

[![Grafana dashboard](https://user-images.githubusercontent.com/1588935/116655813-a7470400-a9c6-11eb-9482-ed6f9369fba2.png)](https://user-images.githubusercontent.com/1588935/116655813-a7470400-a9c6-11eb-9482-ed6f9369fba2.png)

Running on Kubernetes
---

There are example manifests in [k8s](https://github.com/rinx/alvd/tree/main/k8s) directory.

```sh
$ # create new namespace
$ kubectl create ns alvd
$ # change current namespace
$ kubectl config set-context $(kubectl config current-context) --namespace=alvd
$ # deploy Servers
$ kubectl apply -f k8s/server.yaml

$ # after Servers become ready, deploy Agents
$ kubectl apply -f k8s/agent.yaml
```

Lua based config
---

Instead of using command-line flags, alvd can be configured by using a Lua based config file.
Thare's an example Lua file at [examples/config/config.lua](examples/config/config.lua).

```sh
$ ./alvd server --config=examples/config/config.lua
```

### Interceptor features

alvd has interceptor features (filtering, sorting, translating, etc...) that is extensible by using Lua scripts.  
To enable them, run alvd server by passing a path to the Lua scripts.

    $ ./alvd server --config=examples/interceptors/sort.lua

There're various types of examples of interceptors are available in [examples/interceptors](examples/interceptors) directory and [examples/config/config.lua](examples/config/config.lua).

This feature is powered by [yuin/gopher-lua](https://github.com/yuin/gopher-lua) and [vadv/gopher-lua-libs](https://github.com/vadv/gopher-lua-libs).

Current Status
---

- Agent uses NGT service package of Vald Agent NGT.
- uses Vald v1 API on master branch (https://github.com/vdaas/vald/tree/master/apis/proto/v1).
    - Server has APIs in https://github.com/vdaas/vald/tree/master/apis/proto/v1/vald
        - Unary APIs and Streaming APIs are supported.
        - MultiXXX APIs are not supported.
    - Agent has APIs in https://github.com/vdaas/vald/tree/master/apis/proto/v1/vald and https://github.com/vdaas/vald/tree/master/apis/proto/v1/agent/core.
        - Unary APIs and Streaming APIs are supported.
        - MultiXXX APIs are not supported.
- Currently, high-availability (HA) cluster is partly supported.
    - Only static IP multiple Servers (+ multiple Agents) can build an alvd HA cluster.
        - Agents cannot resolve dynamic IPs of Servers.
    - In the near future, HA cluster will be fully supported.


Build
---

    $ make cmd/alvd/alvd


License
---

Same as Vald, alvd is distributed under Apache 2.0 license. (Partially distributed under Mozilla Public License 2.0)

alvd depends on Vald codebase, the files came from Vald (such as `internal`, `pkg/vald`. They are downloaded when running `make` command.) are excluded from my license and ownership.

This is not an official project of Vald. This project is an artifact of 20% project of Vald team.

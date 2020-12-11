alvd - A Lightweight Vald
===

[![License: Apache 2.0](https://img.shields.io/github/license/rinx/alvd.svg?style=flat-square)](https://opensource.org/licenses/Apache-2.0)
[![release](https://img.shields.io/github/release/rinx/alvd.svg?style=flat-square)](https://github.com/rinx/alvd/releases/latest)
[![ghcr.io](https://img.shields.io/badge/ghcr.io-rinx%2Falvd-brightgreen?logo=docker&style=flat-square)](https://github.com/users/rinx/packages/container/package/alvd)
[![Docker Pulls](https://img.shields.io/docker/pulls/rinx/alvd.svg?style=flat-square)](https://hub.docker.com/r/rinx/alvd)

A lightweight distributed vector search engine based on [Vald](https://vald.vdaas.org) codebase.

- works without Kubernetes
- single binary (less than 40MB)
- easy to run (can be configured by command-line options)
- consists of Agent and Server
    - alvd has same features that Vald's gateway-lb + discoverer and agent-ngt have.

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

Quick Start
---

1. Get a latest build from [Actions](https://github.com/rinx/alvd/actions) build results and unzip it.
2. Run alvd Server.
    ```sh
    $ ./alvd server
    2020-12-04 17:30:27     [INFO]: start alvd server
    2020-12-04 17:30:27     [INFO]: websocket server starting on 0.0.0.0:8000
    2020-12-04 17:30:27     [INFO]: start alvd agent
    2020-12-04 17:30:27     [INFO]: gateway gRPC API starting on 0.0.0.0:8080
    2020-12-04 17:30:27     [INFO]: executing daemon pre-start function
    2020-12-04 17:30:27     [INFO]: executing daemon start function
    2020-12-04 17:30:27     [INFO]: server grpc executing preStartFunc
    2020-12-04 17:30:27     [INFO]: gRPC server grpc starting on 0.0.0.0:8081
    INFO[0000] Connecting to proxy                           url="ws://0.0.0.0:8000/connect"
    INFO[0000] Handling backend connection request [e6pv4sgbv4v78soeosb0]
    2020-12-04 17:30:27     [INFO]: connected to: 0.0.0.0:8000
    ```
    alvd Server's websocket server starts on 0.0.0.0:8000 and alvd Server's gRPC API starts on 0.0.0.0:8080.
    Also, alvd Agent's gRPC API starts on 0.0.0.0:8081 (alvd Agent process on the Server can be disabled using `--agent=false` option).
3. Run alvd Agent on a different node (or a different terminal on the same node with `--server 0.0.0.0:8000` and `--grpc-port 8082` option).
    ```sh
    $ ./alvd agent --server host-of-server-node:8000
    $ # ./alvd agent --server 0.0.0.0:8000 --grpc-port 8082
    2020-12-04 17:31:34     [INFO]: start alvd agent
    2020-12-04 17:31:34     [INFO]: executing daemon pre-start function
    2020-12-04 17:31:34     [INFO]: executing daemon start function
    2020-12-04 17:31:34     [INFO]: server grpc executing preStartFunc
    2020-12-04 17:31:34     [INFO]: gRPC server grpc starting on 0.0.0.0:8081
    INFO[0000] Connecting to proxy                           url="ws://host-of-server-node:8000/connect"
    2020-12-04 17:31:34     [INFO]: connected to: host-of-server-node:8000
    ```
4. Add more alvd Agents on the other nodes (or the other ports on the same node).
    ```sh
    $ ./alvd agent --server host-of-server-node:8000
    $ # ./alvd agent --server 0.0.0.0:8000 --grpc-port 808{3,4,5}
    ```
5. Now we can access the alvd Server's gRPC API (`host-of-server-node:8080`) using Vald v1 clients.
    If you don't have one, you can use [valdcli-v1-alpha](https://github.com/vdaas/vald-client-clj/pull/14#issuecomment-738521578) (this CLI is built for linux-amd64).
    ```sh
    $ # insert 100 vectors (dimension: 784) with random IDs
    $ ./valdcli rand-vecs -d 784 -n 100 --with-ids | ./valdcli -h host-of-server-node -p 8080 stream-insert
    $ # search a random vector
    $ ./valdcli rand-vec -d 784 | ./valdcli -h host-of-server-node -p 8080 search
    ```

Current Status
---

- Agent is just wrapping Vald Agent NGT.
- uses Vald v1 API ([#826](https://github.com/vdaas/vald/pull/826)) scheme
    - Server has APIs in https://github.com/vdaas/vald/tree/feature/apis/v1-new-design/apis/proto/v1/vald
        - Unary APIs and Streaming APIs are supported.
        - MultiXXX APIs are not supported.


Build
---

Just running

    $ make cmd/alvd/alvd


License
---

alvd is distributed under Apache 2.0 license, same as Vald.

alvd depends on Vald codebase, the files came from Vald (such as `internal`, `pkg/vald`. They are downloaded when running `make` command.) are excluded from my license and ownership.

This is not an official project of Vald. This project is an artifact of 20% project of Vald team.

+++
title = "Quick Start"
weight = 20
+++

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

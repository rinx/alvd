+++
title = "Docker Compose"
weight = 10
+++

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

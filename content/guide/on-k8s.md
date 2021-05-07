+++
title = "Kubernetes"
weight = 20
+++

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

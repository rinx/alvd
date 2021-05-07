+++
title = "Rationale"
weight = 5
+++

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

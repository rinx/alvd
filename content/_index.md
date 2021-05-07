+++
title = "alvd - A Lightweight Vald"
insert_anchor_links = "right"
+++

<div align="center">
<img src="https://raw.githubusercontent.com/rinx/alvd/main/assets/cover.svg" width="100%">
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

License
---

Same as Vald, alvd is distributed under Apache 2.0 license. (Partially distributed under Mozilla Public License 2.0)

alvd depends on Vald codebase, the files came from Vald (such as `internal`, `pkg/vald`. They are downloaded when running `make` command.) are excluded from my license and ownership.

This is not an official project of Vald. This project is an artifact of 20% project of Vald team.

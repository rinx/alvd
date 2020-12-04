alvd - A Lightweight Vald
===

 [![ghcr.io](https://img.shields.io/badge/ghcr.io-rinx%2Falvd-brightgreen?logo=docker&style=flat-square)](https://github.com/users/rinx/packages/container/package/alvd)

A lightweight distributed vector search engine based on [Vald](https://vald.vdaas.org) codebase.

- single binary
- easy to run
- consists of Agent and Server
- works without Kubernetes

Quick Start
---

Get a latest build from [Actions](https://github.com/rinx/alvd/actions) build results and unzip it.


TBW


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

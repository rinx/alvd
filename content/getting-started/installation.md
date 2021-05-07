+++
title = "Installation"
weight = 30
+++

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

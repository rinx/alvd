+++
title = "Egress filter"
weight = 50
+++

Lua based config
---

Instead of using command-line flags, alvd can be configured by using a Lua based config file.
Thare's an example Lua file at [examples/config/config.lua](https://github.com/rinx/alvd/tree/main/examples/config/config.lua).

    $ ./alvd server --config=examples/config/config.lua

### Egress filter feature

alvd has an egress filter (= post filter) feature (filtering, sorting, translating, etc...) that is extensible by using Lua scripts.

To enable it, run alvd server by passing a path to the Lua scripts.

    $ ./alvd server --egress-filter-lua-filepath=examples/egress-filter/sort.lua

There're various types of examples of filters are available in [examples/egress-filter](https://github.com/rinx/alvd/tree/main/examples/egress-filter) directory.

This feature is powered by [yuin/gopher-lua](https://github.com/yuin/gopher-lua) and [vadv/gopher-lua-libs](https://github.com/vadv/gopher-lua-libs).

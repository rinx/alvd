+++
title = "Lua based config"
weight = 50
+++

Lua based config
---

Instead of using command-line flags, alvd can be configured by using a Lua based config file.
Thare's an example Lua file at [examples/config/config.lua](https://github.com/rinx/alvd/tree/main/examples/config/config.lua).

    $ ./alvd server --config=examples/config/config.lua

### Interceptor features

alvd has interceptor features (filtering, sorting, translating, etc...) that is extensible by using Lua scripts.  
To enable them, run alvd server by passing a path to the Lua scripts.

    $ ./alvd server --config=examples/interceptors/sort.lua

There're various types of examples of interceptors are available in [examples/interceptors](https://github.com/rinx/alvd/tree/main/examples/interceptors) directory and [examples/config/config.lua](https://github.com/rinx/alvd/tree/main/examples/config/config.lua).

This feature is powered by [yuin/gopher-lua](https://github.com/yuin/gopher-lua) and [vadv/gopher-lua-libs](https://github.com/vadv/gopher-lua-libs).


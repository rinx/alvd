-- modules in gopher-lua-libs are available in Lua scripts.
-- https://github.com/vadv/gopher-lua-libs
local json = require("json")
local time = require("time")

agent = {
  server_addresses = {"0.0.0.0:8000"},
  agent_name = "",
  log_level = "info",
  dimension = 300,
  distance_type = "l2",
  object_type = "float",
  creation_edge_size = 10,
  search_edge_size = 20,
  bulk_insert_chunk_size = 100,
  index_path = "",
  index_self_check_interval = "30m",
  grpc_host = "0.0.0.0",
  grpc_port = 8081,
  metrics_host = "0.0.0.0",
  metrics_port = 9090,
  metrics_collect_interval = "5s",
}

server = {
  agent_enabled = true,
  log_level = "info",

  server_grpc_host = "0.0.0.0",
  server_grpc_port = 8080,

  metrics_host = "0.0.0.0",
  metrics_port = 9090,
  metrics_collect_interval = "5s",

  replicas = 2,
  check_index_interval = "10s",
  create_index_threshold = 1000,

  -- server-side Search Result Interceptor
  -- it can be used for post-filtering, sorting,
  -- translating or modifying search results.
  search_result_interceptor = function (config, results, retry)
    -- config: search config
    -- results: returned search results
    -- retry: retry config
    for i, r in results() do
      results[i].Id = json.encode(
        {
          id = r.Id,
          time = time.format(
            time.unix(),
            "Jan  2 15:04:05 2006",
            "Asia/Tokyo"),
        }
      )
    end
  end,

  -- server-side Search Query Interceptor
  search_query_interceptor = function (request)
    -- print(string.format("Searching top %d neighbors", request.Config.Num))
  end,

  -- server-side Insert Data Interceptor
  insert_data_interceptor = function (request)
    -- print(string.format("Inserting ID: %s", request.Vector.Id))
  end,
}

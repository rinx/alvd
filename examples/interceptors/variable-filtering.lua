local json = require("json")

server = {
  -- this interceptor filters out results that has larger distance than
  -- the specified threshold.
  -- threshold can be specified in JSON formatted 'EgressFilters' field
  -- in the search requests.
  --
  -- eg.)
  -- Search_Config.EgressFilters.Targets[0].Host = "{\"threshold\": 10.8}"
  search_result_interceptor = function (config, results, retry)
    if config.EgressFilters and
      config.EgressFilters.Targets and
      #config.EgressFilters.Targets > 0 then
      local cfg, err = json.decode(config.EgressFilters.Targets[1].Host)
      if err then error(err) end

      print(string.format("executing distance-filtering by threshold %f", cfg.threshold))

      local remains = {}
      for i, r in results() do
        -- remove elements by distances
        if r.Distance < cfg.threshold then
          remains[#remains+1] = r
        end

        results[i] = nil
      end

      for i, r in pairs(remains) do
        results[i] = r
      end
    end
  end,
}

server = {
  search_result_interceptor = function (config, results, retry)
    -- if egress filter is set, execute filtering
    if config.EgressFilters and
      config.EgressFilters.Targets and
      #config.EgressFilters.Targets > 0 then
      if config.EgressFilters.Targets[1].Host == "distance-filtering" then
        print("executing distance-filtering")
        local remains = {}
        for i, r in results() do
          -- remove elements by distances
          if r.Distance < 10.9 then
            remains[#remains+1] = r
          end

          results[i] = nil
        end

        for i, r in pairs(remains) do
          results[i] = r
        end
      end
    end
  end,
}

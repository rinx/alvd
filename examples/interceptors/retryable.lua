server = {
  search_result_interceptor = function (config, results, retry)
    -- if `retry.Enabled` is true, retry ANN search when the number of `results` is lower than the required number.
    retry.Enabled = true
    -- `retry.MaxRetries` represents maximum number of retries.
    retry.MaxRetries = 3
    -- `retry.NextNumMultiplier` represents how to increase number of internal search results.
    retry.NextNumMultiplier = 2

    local remains = {}
    for i, r in results() do
      -- remove elements if ID lengths is lower than 3
      if string.len(r.Id) >= 3 then
        remains[#remains+1] = r
      end

      results[i] = nil
    end

    for i, r in pairs(remains) do
      results[i] = r
    end
  end,
}

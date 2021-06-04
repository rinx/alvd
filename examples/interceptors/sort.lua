server = {
  search_result_interceptor = function (results, retry)
    local sorter = {}
    for i, r in results() do
      sorter[i] = r
    end

    -- sort
    table.sort(sorter, function(a, b)
      -- reverse order
      return a.Distance > b.Distance
    end)

    -- put the sorted data into `results`
    for i, r in pairs(sorter) do
      results[i] = r
    end
  end,
}

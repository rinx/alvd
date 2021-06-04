server = {
  search_result_interceptor = function (results, retry)
    for i, r in results() do
      results[i].Id = string.reverse(r.Id)
    end
  end,
}

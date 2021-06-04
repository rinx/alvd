server = {
  search_result_interceptor = function (results, retry)
    for i, r in results() do
      print(string.format("Id: %s, Distance: %3.7f", r.Id, r.Distance))
    end
  end,
}

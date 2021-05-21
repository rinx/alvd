server = {
  egress_filter = function (results, retry)
    for i, r in results() do
      results[i].Id = string.reverse(r.Id)
    end
  end,
}

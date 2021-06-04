server = {
  search_result_interceptor = function (config, results, retry)
    for i, r in results() do
      print(string.format("Id: %s, Distance: %3.7f", r.Id, r.Distance))
    end
  end,

  search_query_interceptor = function (request)
    print(string.format("Searching top %d neighbors", request.Config.Num))
  end,

  insert_data_interceptor = function (request)
    print(string.format("Inserting ID: %s", request.Vector.Id))
  end,

}

-- using vadv/gopher-lua-libs
local json = require("json")
local time = require("time")

server = {
  search_result_interceptor = function (config, results, retry)
    for i, r in results() do
      results[i].Id = json.encode(
        {
          id = r.Id,
          time = time.format(time.unix(), "Jan  2 15:04:05 2006", "Asia/Tokyo")
        }
      )
    end
  end,
}

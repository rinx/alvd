-- `results` is a slice of search results.

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

-- dump.lua is useful for debugging.
-- `results` is a slice of search results.

for i, r in results() do
  print(string.format("Id: %s, Distance: %3.7f", r.Id, r.Distance))
end

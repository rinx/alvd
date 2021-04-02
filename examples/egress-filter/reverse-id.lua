-- `results` is a slice of search results.

for i, r in results() do
  results[i].Id = string.reverse(r.Id)
end

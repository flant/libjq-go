def camel2:
  gsub("-(?<a>[a-z])"; .a|ascii_upcase);

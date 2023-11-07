package main

var commentMap = map[string][]string{
  "javascript":     {"//", "", "/*", "*/", "", ""},
  "python":         {"#", "", `"""`, `"""`, "'''", "'''"},
  "java":           {"//", "", "/*", "*/", "", ""},
  "c":              {"//", "", "/*", "*/", "", ""},
  "c++":            {"//", "", "/*", "*/", "", ""},
  "c#":             {"//", "", "/*", "*/", "", ""},
  "php":            {"//", "#", "/*", "*/", "", ""},
  "typescript":     {"//", "", "/*", "*/", "", ""},
  "go":             {"//", "", "/*", "*/", "", ""},
  "ruby":           {"#", "", "=begin", "=end", "", ""},
  "swift":          {"//", "", "/*", "*/", "", ""},
  "rust":           {"//", "", "/*", "*/", "", ""},
  "kotlin":         {"//", "", "/*", "*/", "", ""},
  "dart":           {"//", "", "/*", "*/", "", ""},
  "r":              {"#", "", "", "", "", ""},
  "shell":          {"#", "", "", "", "", ""},
  "matlab":         {"%", "", "%{", "%}", "", ""},
  "vba":            {"'", "", "", "", "", ""},
  "scala":          {"//", "", "/*", "*/", "", ""},
  "perl":           {"#", "", "=pod", "=cut", "", ""},
  "haskell":        {"--", "", "{-", "-}", "", ""},
  "lua":            {"--", "", "--[[", "]]", "", ""},
  "groovy":         {"//", "", "/*", "*/", "", ""},
  "coffeescript":   {"#", "", "###", "###", "", ""},
  "raku":           {"#", "", "=begin", "=end", "", ""},
  "objective-c":    {"//", "", "/*", "*/", "", ""},
  "sql":            {"--", "", "/*", "*/", "", ""},
  "powershell":     {"#", "", "<#", "#>", "", ""},
  "julia":          {"#", "", "#=", "=#", "", ""},
  "fortran":        {"!", "", "", "", "", ""},
  "ada":            {"--", "", "", "", "", ""},
  "elixir":         {"#", "", "", "", "", ""},
  "erlang":         {"%", "", "", "", "", ""},
  "f#":             {"//", "", "(*", "*)", "", ""},
  "crystal":        {"#", "", "", "", "", ""},
  "apex":           {"//", "", "/*", "*/", "", ""},
  "ocaml":          {"(*", "", "(*", "*)", "", ""},
  "ballerina":      {"//", "", "/*", "*/", "", ""},
  "nim":            {"#", "", "", "", "", ""},
  "d":              {"//", "", "/*", "*/", "", ""},
  "clojure":        {";", "", "", "", "", ""},
  "pascal":         {"//", "", "{", "}", "", ""},
  "delphi":         {"//", "", "{", "}", "", ""},
  "prolog":         {"%", "", "/*", "*/", "", ""},
  "elm":            {"--", "", "{-", "-}", "", ""},
  "scheme":         {";", "", "#|", "|#", "", ""},
  "lisp":           {";", "", "#|", "|#", "", ""},
  "vb.net":         {"'", "", "", "", "", ""},
  "bash":           {"#", "", "", "", "", ""},
  "html":           {"", "", "<!--", "-->", "", ""},
  "css":            {"", "", "/*", "*/", "", ""},
  "xml":            {"", "", "<!--", "-->", "", ""},
  "json":           {"", "", "", "", "", ""},
  "yaml":           {"#", "", "", "", "", ""},
  "markdown":       {"", "", "<!--", "-->", "", ""},
  "latex":          {"%", "", "", "", "", ""},
  "sass":           {"//", "", "/*", "*/", "", ""},
  "scss":           {"//", "", "/*", "*/", "", ""},
  "less":           {"//", "", "/*", "*/", "", ""},
  "stylus":         {"//", "", "/*", "*/", "", ""},
  "assembly":       {";", "", "", "", "", ""},
  "autoit":         {";", "", "", "", "", ""},
  "batch":          {"REM", "::", "", "", "", ""},
  "toml":           {"#", "", "", "", "", ""},
  "ini":            {";", "#", "", "", "", ""},
  "dockerfile":     {"#", "", "", "", "", ""},
  "makefile":       {"#", "", "", "", "", ""},
  "terraform":      {"#", "", "/*", "*/", "", ""},
  "ansible":        {"#", "", "", "", "", ""},
}
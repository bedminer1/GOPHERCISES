module example.com/main

go 1.22.1

replace example.com/parse => ../../link-parser/parse

require example.com/parse v0.0.0-00010101000000-000000000000

require golang.org/x/net v0.22.0 // indirect

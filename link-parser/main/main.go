package main

import (
	"fmt"
	"strings"

	link "example.com/parse"
)

var exampleHTML = `
<html>
	<body>
		<a href="github.com/bedminer1">a link</a>
	</body>
</html>
`

func main() {
	r := strings.NewReader(exampleHTML)
	links, err := link.Parse(r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", links)
}
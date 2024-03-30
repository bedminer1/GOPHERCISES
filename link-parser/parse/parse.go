package link

import (
	"io"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

// take HTML doc and return slice of links
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	// filter out everything with <a> tag in html
	nodes := linkNodes(doc)
	var links []Link

	for _, node := range nodes {
		links = append(links, buildLink(node))
	}
	
	return links, nil
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}

func buildLink(n *html.Node) Link {
	var ret Link
	// parse href
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}
	ret.Text = parseText(n)
	return ret
}

func parseText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += parseText(c) + " "
	}
	return ret
}
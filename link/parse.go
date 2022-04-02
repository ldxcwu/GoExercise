package link

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

//Link represents a link (<a href="...">) in an HTML document
type Link struct {
	Href string
	Text string
}

//Parse will take in an HTML document and will return a slice
//of links parsed from it
func Parse(r io.Reader) ([]Link, error) {
	node, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := linkNodes(node)
	var links []Link
	for _, n := range nodes {
		links = append(links, buildLink(n))
	}
	fmt.Println(links)
	return nil, nil
}

//解析标签树，存于切片
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

//将Node映射为Link，之前解析的Node全部都是a标签，
//但是也只是解析到a标签的程度，
func buildLink(n *html.Node) Link {
	var ret Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}
	ret.Text = text(n)
	return ret
}

func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	//a标签内部并未进行递归解析
	// var ret string
	var buf bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// ret += text(c) + " "
		buf.WriteString(text(c))
	}
	//strings.Fields(s)会将s按照空白符(1个或多个)进行分割，
	return strings.Join(strings.Fields(buf.String()), " ")
}

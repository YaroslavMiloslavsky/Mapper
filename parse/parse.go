package parse

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type hyperText interface {
	GetHref() (string, error)
	GetText() (string, error)
}

// <a href>
type Link struct {
	href string
	text string
}

/*
	Todo: Write tests!!!!!
*/

// <h#>
type Header struct {
	text string
}

// <div>
// type Div struct {
// }

func (l *Link) GetHref() (string, error) {
	if len(l.href) == 0 {
		return "", fmt.Errorf("there is no reference")
	}
	return l.href, nil
}

func (l *Link) GetText() (string, error) {
	if len(l.text) == 0 {
		return "", fmt.Errorf("there is no text")
	}
	return l.text, nil
}

func (h *Header) GetHref() (string, error) {
	return "", nil
}

func (h *Header) GetText() (string, error) {
	if len(h.text) == 0 {
		return "", fmt.Errorf("there is no text")
	}
	return h.text, nil
}

func ParseHeaders(file *os.File) ([]Header, error) {
	reader := bufio.NewReader(file)
	doc, err := html.Parse(reader)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	nodes := headerNodes(doc)
	headers := []Header{}

	for _, node := range nodes {
		headers = append(headers, buildHeader(node))
	}

	file.Seek(0, 0)
	return headers, nil
}

func headerNodes(node *html.Node) []*html.Node {
	if node.Type == html.ElementNode && (node.Data == "h1" || node.Data == "h2" || node.Data == "h3" || node.Data == "h4") {
		return []*html.Node{node}
	}
	var retVal []*html.Node
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		retVal = append(retVal, headerNodes(child)...)
	}
	return retVal
}

func buildHeader(node *html.Node) Header {
	var text string
	// Text value
	text = getText(node)

	return Header{text: text}
}

func ParseLinks(file *os.File) ([]Link, error) {
	reader := bufio.NewReader(file)
	doc, err := html.Parse(reader)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	nodes := linkNodes(doc)
	link := []Link{}
	for _, node := range nodes {
		link = append(link, buildLink(node))
	}

	file.Seek(0, 0)
	return link, nil
}

func buildLink(node *html.Node) Link {
	var href, text string
	// Href value
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			href = attr.Val
		}
	}
	// Text value
	text = getText(node)

	return Link{href: href, text: text}

}

func linkNodes(node *html.Node) []*html.Node {
	if node.Type == html.ElementNode && node.Data == "a" {
		return []*html.Node{node}
	}
	var retVal []*html.Node
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		retVal = append(retVal, linkNodes(child)...)
	}
	return retVal
}

func getText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}
	var retVal string
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		retVal += strings.TrimSpace(getText(child)) + "  "

	}
	return strings.Join(strings.Fields(retVal), " ")
}

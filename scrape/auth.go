package scrape

import (
	"io"

	"github.com/PuerkitoBio/goquery"
)

// return RelayState, SAMLResponse
func RelayStateAndSAMLResponse(htmlReader io.ReadCloser) (string, string, error) {
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return "", "", err
	}
	selection := doc.Find("html > body > form > div > input")
	relayState, ok := selection.Attr("value")
	if !ok {
		return "", "", &ErrorNotFound{Name: "RelayState"}
	}
	selection = selection.Next()
	samlResponse, ok := selection.Attr("value")
	if !ok {
		return "", "", &ErrorNotFound{Name: "SAMLResponse"}
	}

	return relayState, samlResponse, nil
}

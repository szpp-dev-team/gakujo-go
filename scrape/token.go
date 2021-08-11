package scrape

import (
	"io"

	"github.com/PuerkitoBio/goquery"
)

func ApacheToken(htmlReader io.ReadCloser) (string, error) {
	// ページによってtokenの場所が違う場合
	selectors := []string{
		"#SC_A01_06 > form:nth-child(15) > div > input[type=hidden]",
		"#header > form:nth-child(4) > div > input[type=hidden]",
	}
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return "", err
	}
	for _, selector := range selectors {
		selection := doc.Find(selector)
		token, ok := selection.Attr("value")
		if ok {
			return token, nil
		}
	}
	return "", &ErrorNotFound{Name: "org.apache.struts.taglib.html.TOKEN"}
}

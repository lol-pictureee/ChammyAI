package quotes

import (
	"bot_tg/routes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetQuote() (string, error) {
	resp, err := http.Get(routes.ApiQuotes.String())
	if err != nil {
		return "", fmt.Errorf("HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("code error: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error response body: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return "", fmt.Errorf("error creating document: %w", err)
	}

	var quote string
	doc.Find("div.field-item.even.last").Each(func(i int, s *goquery.Selection) {

		if quote != "" {
			return
		}
		quote = s.Text()
	})

	if quote == "" {
		return "", fmt.Errorf("Проверьте tag")
	}

	return quote, nil
}

package agent

import (
	"fmt"
	"github.com/gocolly/colly"
	"strings"
)

type SearchItem struct {
	Title       string
	URL         string
	Description string
}

type GoogleSearchEngine struct {
}

func (g *GoogleSearchEngine) PerformSearch(query string, numResults int) ([]SearchItem, error) {
	c := colly.NewCollector()

	var results []SearchItem

	c.OnHTML("div.g", func(e *colly.HTMLElement) {
		if len(results) >= numResults {
			return
		}

		title := e.ChildText("h3")
		url := e.ChildAttr("a", "href")
		description := e.ChildText("div.IsZvec")

		if title != "" && url != "" {
			results = append(results, SearchItem{
				Title:       title,
				URL:         url,
				Description: description,
			})
		}
	})

	err := c.Visit(fmt.Sprintf("https://www.google.com/search?q=%s", strings.ReplaceAll(query, " ", "+")))
	if err != nil {
		return nil, err
	}

	return results, nil
}

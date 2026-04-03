package adapter

import (
	"context"
	"time"

	"github.com/mmcdole/gofeed"
)

// RawArticle is the output of a feed adapter.
type RawArticle struct {
	Title     string
	Content   string
	Link      string
	PubDate   time.Time
	Source    string
}

// RSSAdapter fetches articles from an RSS/Atom feed.
type RSSAdapter struct {
	Name    string
	URL     string
	Enabled bool
	parser  *gofeed.Parser
}

func NewRSSAdapter(name, url string, enabled bool) *RSSAdapter {
	return &RSSAdapter{
		Name:    name,
		URL:     url,
		Enabled: enabled,
		parser:  gofeed.NewParser(),
	}
}

func (a *RSSAdapter) IsEnabled() bool {
	return a.Enabled
}

func (a *RSSAdapter) GetName() string {
	return a.Name
}

func (a *RSSAdapter) Fetch(ctx context.Context) ([]*RawArticle, error) {
	if !a.Enabled {
		return nil, nil
	}

	feed, err := a.parser.ParseURLWithContext(a.URL, ctx)
	if err != nil {
		return nil, err
	}

	var articles []*RawArticle
	for _, item := range feed.Items {
		pubDate := time.Now()
		if item.PublishedParsed != nil {
			pubDate = *item.PublishedParsed
		}

		content := item.Description
		if item.Content != "" {
			content = item.Content
		}

		articles = append(articles, &RawArticle{
			Title:   item.Title,
			Content: content,
			Link:    item.Link,
			PubDate: pubDate,
			Source:  a.Name,
		})
	}

	return articles, nil
}

package crawler

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/lunancy1992/jianghu-server/internal/ai"
	"github.com/lunancy1992/jianghu-server/internal/crawler/adapter"
	"github.com/lunancy1992/jianghu-server/internal/model"
	"github.com/lunancy1992/jianghu-server/internal/repo"
)

// Adapter is the interface for feed adapters.
type Adapter interface {
	IsEnabled() bool
	GetName() string
	Fetch(ctx context.Context) ([]*adapter.RawArticle, error)
}

// Crawler fetches, deduplicates, and publishes news.
type Crawler struct {
	adapters []Adapter
	newsRepo *repo.NewsRepo
	aiClient ai.AIClient
}

func NewCrawler(newsRepo *repo.NewsRepo, aiClient ai.AIClient) *Crawler {
	return &Crawler{
		newsRepo: newsRepo,
		aiClient: aiClient,
	}
}

func (c *Crawler) AddAdapter(a Adapter) {
	c.adapters = append(c.adapters, a)
}

// Run executes one crawl cycle across all adapters.
func (c *Crawler) Run(ctx context.Context) {
	log.Println("[Crawler] Starting crawl cycle")

	var totalNew, totalSkipped, totalErrors int

	for _, a := range c.adapters {
		if !a.IsEnabled() {
			continue
		}

		log.Printf("[Crawler] Fetching from %s", a.GetName())
		articles, err := a.Fetch(ctx)
		if err != nil {
			log.Printf("[Crawler] Error fetching %s: %v", a.GetName(), err)
			totalErrors++
			continue
		}

		for _, raw := range articles {
			fp := fmt.Sprintf("%016x", SimHash(raw.Title+raw.Content))

			// Check dedup
			existing, err := c.newsRepo.FindByFingerprint(ctx, fp)
			if err != nil {
				log.Printf("[Crawler] Dedup check error: %v", err)
				totalErrors++
				continue
			}
			if existing != nil {
				totalSkipped++
				continue
			}

			// AI summarize
			news := &model.News{
				Title:       raw.Title,
				Content:     cleanHTML(raw.Content),
				Source:      raw.Source,
				SourceURL:   raw.Link,
				Fingerprint: fp,
				Status:      1, // published
			}

			if c.aiClient != nil {
				summary, err := c.aiClient.Summarize(ctx, raw.Title, news.Content)
				if err != nil {
					log.Printf("[Crawler] AI summarize error: %v", err)
					news.Summary = raw.Title
					news.Category = "其他"
				} else {
					news.Summary = summary.OneLiner
					news.Category = summary.Category
					news.Keywords = toJSONArray(summary.Keywords)
					news.Sensitivity = summary.Sensitivity
				}
			} else {
				news.Summary = raw.Title
				news.Category = "其他"
			}

			pubTime := raw.PubDate
			news.PublishedAt = &pubTime

			_, err = c.newsRepo.Create(ctx, news)
			if err != nil {
				log.Printf("[Crawler] Save error: %v", err)
				totalErrors++
				continue
			}
			totalNew++
		}
	}

	log.Printf("[Crawler] Cycle complete: %d new, %d skipped, %d errors", totalNew, totalSkipped, totalErrors)
}

func cleanHTML(s string) string {
	// Basic HTML tag removal - for production use a real HTML sanitizer
	result := s
	for {
		start := strings.Index(result, "<")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], ">")
		if end == -1 {
			break
		}
		result = result[:start] + result[start+end+1:]
	}
	return strings.TrimSpace(result)
}

func toJSONArray(items []string) string {
	if len(items) == 0 {
		return "[]"
	}
	parts := make([]string, len(items))
	for i, item := range items {
		parts[i] = `"` + strings.ReplaceAll(item, `"`, `\"`) + `"`
	}
	return "[" + strings.Join(parts, ",") + "]"
}

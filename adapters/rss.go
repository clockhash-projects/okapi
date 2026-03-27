package adapters

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"okapi/internal/models"
)

type RSSConfig struct {
	ID                  string `yaml:"id"`
	DisplayName         string `yaml:"display_name"`
	Kind                string `yaml:"kind"`
	URL                 string `yaml:"url"`
	PollIntervalSeconds int    `yaml:"poll_interval_seconds"`
}

type RSSAdapter struct {
	cfg        RSSConfig
	httpClient *http.Client
}

func NewRSSAdapter(cfg RSSConfig) *RSSAdapter {
	return &RSSAdapter{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (a *RSSAdapter) ID() string          { return a.cfg.ID }
func (a *RSSAdapter) DisplayName() string { return a.cfg.DisplayName }
func (a *RSSAdapter) PollInterval() time.Duration {
	return time.Duration(a.cfg.PollIntervalSeconds) * time.Second
}

type rssFeed struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Items []struct {
			Title       string `xml:"title"`
			Description string `xml:"description"`
			PubDate     string `xml:"pubDate"`
			GUID        string `xml:"guid"`
		} `xml:"item"`
	} `xml:"channel"`
}

type atomFeed struct {
	XMLName  xml.Name `xml:"feed"`
	Title    string   `xml:"title"`
	Subtitle string   `xml:"subtitle"`
	Updated  string   `xml:"updated"`
	Entries  []struct {
		ID      string `xml:"id"`
		Title   string `xml:"title"`
		Updated string `xml:"updated"`
		Content string `xml:"content"`
		Summary string `xml:"summary"`
	} `xml:"entry"`
}

func (a *RSSAdapter) Fetch(ctx context.Context) (*models.StatusResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", a.cfg.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upstream returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	res := &models.StatusResponse{
		Service:    a.cfg.ID,
		Status:     models.StatusOperational,
		Summary:    "All systems operational",
		FetchedAt:  time.Now(),
		DataSource: "feed",
		SourceURL:  a.cfg.URL,
	}

	// Try RSS first
	var rss rssFeed
	if err := xml.Unmarshal(body, &rss); err == nil && rss.XMLName.Local == "rss" {
		if len(rss.Channel.Items) > 0 {
			lastItem := rss.Channel.Items[0]
			pubDate, _ := time.Parse(time.RFC1123, lastItem.PubDate)
			if pubDate.IsZero() {
				// Try some other formats common in RSS
				pubDate, _ = time.Parse(time.RFC822, lastItem.PubDate)
			}

			if !pubDate.IsZero() && time.Since(pubDate) < 24*time.Hour {
				res.Status = models.StatusDegraded
				res.Summary = lastItem.Title
			}

			for _, item := range rss.Channel.Items {
				createdAt, _ := time.Parse(time.RFC1123, item.PubDate)
				res.Incidents = append(res.Incidents, models.Incident{
					ID:        item.GUID,
					Title:     item.Title,
					Status:    "active",
					Body:      item.Description,
					CreatedAt: createdAt,
				})
			}
		}
		return res, nil
	}

	// Try Atom
	var atom atomFeed
	if err := xml.Unmarshal(body, &atom); err == nil && atom.XMLName.Local == "feed" {
		if len(atom.Entries) > 0 {
			lastEntry := atom.Entries[0]
			updated, _ := time.Parse(time.RFC3339, lastEntry.Updated)
			if !updated.IsZero() && time.Since(updated) < 24*time.Hour {
				res.Status = models.StatusDegraded
				res.Summary = lastEntry.Title
			}

			for _, entry := range atom.Entries {
				createdAt, _ := time.Parse(time.RFC3339, entry.Updated)
				body := entry.Summary
				if body == "" {
					body = entry.Content
				}
				res.Incidents = append(res.Incidents, models.Incident{
					ID:        entry.ID,
					Title:     entry.Title,
					Status:    "active",
					Body:      body,
					CreatedAt: createdAt,
				})
			}
		}
		return res, nil
	}

	return nil, fmt.Errorf("failed to parse as RSS or Atom")
}

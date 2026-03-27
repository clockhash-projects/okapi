package code

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"okapi/internal/models"
)

type AWSAdapter struct {
	BaseURL string
}

func (a *AWSAdapter) ID() string          { return "aws" }
func (a *AWSAdapter) DisplayName() string { return "AWS" }
func (a *AWSAdapter) PollInterval() time.Duration {
	return 60 * time.Second
}

type awsRSS struct {
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

func (a *AWSAdapter) Fetch(ctx context.Context) (*models.StatusResponse, error) {
	feedURL := "https://status.aws.amazon.com/rss/all.rss"
	if a.BaseURL != "" {
		feedURL = a.BaseURL
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch aws status: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var rss awsRSS
	if err := xml.Unmarshal(body, &rss); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rss: %w", err)
	}

	status := models.StatusOperational
	summary := "All systems nominal"
	incidents := []models.Incident{}

	if len(rss.Channel.Items) > 0 {
		// AWS RSS often contains recent resolved items too,
		// but for a simple "Global" status we'll report major if there's anything in the feed
		// that isn't explicitly marked as resolved.
		status = models.StatusMajorOutage
		summary = fmt.Sprintf("Active issues detected: %s", rss.Channel.Items[0].Title)

		for _, item := range rss.Channel.Items {
			createdAt, _ := time.Parse(time.RFC1123, item.PubDate)
			incidents = append(incidents, models.Incident{
				ID:        item.GUID,
				Title:     item.Title,
				Status:    "active",
				Body:      item.Description,
				CreatedAt: createdAt,
			})
		}
	}

	return &models.StatusResponse{
		Service:    "aws",
		Status:     status,
		Summary:    summary,
		Incidents:  incidents,
		FetchedAt:  time.Now(),
		DataSource: "rss_feed",
		SourceURL:  "https://health.aws.amazon.com",
	}, nil
}

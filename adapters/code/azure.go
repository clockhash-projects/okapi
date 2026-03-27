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

type AzureAdapter struct {
	BaseURL string
}

func (a *AzureAdapter) ID() string          { return "azure" }
func (a *AzureAdapter) DisplayName() string { return "Azure" }
func (a *AzureAdapter) PollInterval() time.Duration {
	return 60 * time.Second
}

type azureRSS struct {
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

func (a *AzureAdapter) Fetch(ctx context.Context) (*models.StatusResponse, error) {
	feedURL := "https://azure.status.microsoft/en-us/status/feed/"
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
		return nil, fmt.Errorf("failed to fetch azure status: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var rss azureRSS
	if err := xml.Unmarshal(body, &rss); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rss: %w", err)
	}

	status := models.StatusOperational
	summary := "All systems operational"
	incidents := []models.Incident{}

	if len(rss.Channel.Items) > 0 {
		status = models.StatusMajorOutage // Azure RSS feed usually only has items during issues
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
		Service:    "azure",
		Status:     status,
		Summary:    summary,
		Incidents:  incidents,
		FetchedAt:  time.Now(),
		DataSource: "rss_feed",
		SourceURL:  "https://azure.status.microsoft/en-us/status",
	}, nil
}

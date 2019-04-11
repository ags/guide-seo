package guide

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	apiKey string
	http   *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey, &http.Client{
		Timeout: 2500 * time.Millisecond,
	}}
}

type Collection struct {
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	Destinations []Destination `json:"destinations"`
}

type Destination struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Description  string   `json:"description"`
	Website      string   `json:"webSite"`
	BannerImages []string `json:"bannerImages"`
	Latitude     float64  `json:"latitude"`
	Longitude    float64  `json:"longitude"`
	Street       string   `json:"street"`
	Suburb       string   `json:"suburb"`
	State        string   `json:"state"`
	PostCode     string   `json:"postcode"`
}

type FindCollectionInput struct {
	CompanyAPIKey string
	RegionID      int
	CollectionID  int
}

func (c *Client) FindCollection(ctx context.Context, in FindCollectionInput) (Collection, error) {
	url := fmt.Sprintf(
		"https://guide.app/api/v1/regions/%d/collections/%d?type=Collection",
		in.RegionID, in.CollectionID,
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Collection{}, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", c.apiKey)
	req.Header.Set("companyKey", in.CompanyAPIKey)

	res, err := c.http.Do(req)
	if err != nil {
		return Collection{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var e Error
		_ = json.NewDecoder(res.Body).Decode(&e)
		return Collection{}, e
	}

	var co Collection
	err = json.NewDecoder(res.Body).Decode(&co)
	return co, err
}

type Error struct {
	ErrorMessage string   `json:"errorMessage"`
	Errors       []string `json:"errors"`
}

func (e Error) Error() string {
	return e.ErrorMessage
}

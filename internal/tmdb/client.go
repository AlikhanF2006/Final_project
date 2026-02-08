package tmdb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

func (c *Client) GetPopularMovies() ([]MovieDTO, error) {
	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/movie/popular?api_key=%s&language=en-US&page=1",
		c.apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Results []MovieDTO `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Results, nil
}

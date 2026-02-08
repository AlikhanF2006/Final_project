package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/AlikhanF2006/Final_project/internal/postgres/dto"
)

var ErrTMDBRequestFailed = errors.New("tmdb request failed")

type TMDBClient struct {
	token string
}

func NewTMDBClient(token string) *TMDBClient {
	return &TMDBClient{token: token}
}

func (c *TMDBClient) doRequest(url string, target any) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrTMDBRequestFailed
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func (c *TMDBClient) GetMovie(tmdbID int) (dto.TMDBMovieResponse, error) {
	var movie dto.TMDBMovieResponse

	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/movie/%d?language=en-US",
		tmdbID,
	)

	if err := c.doRequest(url, &movie); err != nil {
		return dto.TMDBMovieResponse{}, err
	}

	return movie, nil
}

func (c *TMDBClient) GetTrailerKey(tmdbID int) (string, error) {
	var videos dto.TMDBVideosResponse

	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/movie/%d/videos",
		tmdbID,
	)

	if err := c.doRequest(url, &videos); err != nil {
		return "", err
	}

	for _, v := range videos.Results {
		if v.Site == "YouTube" && v.Type == "Trailer" {
			return v.Key, nil
		}
	}

	return "", nil
}

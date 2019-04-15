package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	es7 "github.com/elastic/go-elasticsearch/v7"
)

type (
	Repository struct {
		es *es7.Client
	}

	responseModel struct {
		Hits ExternalHits `json:"hits"`
	}
	ExternalHits struct {
		Hits []Hits `json:"hits"`
	}
	Hits struct {
		ID     string    `json:"_id"`
		Score  float64   `json:"_score"`
		Source UserModel `json:"_source"`
	}
	UserModel struct {
		Login string `json:"login"`
		Pass  string `json:"pass"`
	}
)

func NewRepository(es *es7.Client) *Repository {
	return &Repository{es: es}
}

func (r *Repository) GetID(ctx context.Context, login, pass string) (id string, resErr error) {
	resp, err := r.es.Search(
		r.es.Search.WithContext(context.Background()),
		r.es.Search.WithIndex("users"),
		r.es.Search.WithBody(strings.NewReader(fmt.Sprintf(`{"query": {"bool": {"must":[{"term": {"login": "%s"}},{"term": {"pass": "%s"}}]}}}`, login, pass))),
		r.es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			resErr = err
		}
	}()

	if resp.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return "", fmt.Errorf("error parsing the response body: %s", err)
		} else {
			return "", fmt.Errorf("error from DB: %v", e["error"].(map[string]interface{})["reason"])
		}
	}

	var res responseModel
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("error parsing the user response body: %s", err)
	}

	if len(res.Hits.Hits) == 0 {
		return "", ErrNotFound{login: login}
	}

	if ln := len(res.Hits.Hits); ln > 1 {
		return "", fmt.Errorf("too much results for login %s: %d", login, ln)
	}

	return res.Hits.Hits[0].ID, nil
}

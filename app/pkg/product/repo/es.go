package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	es7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/kainobor/estest/app/pkg/product"
)

type (
	Repository struct {
		es *es7.Client
	}

	listQuery struct {
		From  *int64                 `json:"from,omitempty"`
		Size  *int64                 `json:"size,omitempty"`
		Sort  []map[string]string    `json:"sort,omitempty"`
		Query map[string]interface{} `json:"query"`
	}

	responseModel struct {
		Hits ExternalHits `json:"hits"`
	}
	ExternalHits struct {
		Hits []Hits `json:"hits"`
	}
	Hits struct {
		ID     string        `json:"_id"`
		Score  float64       `json:"_score"`
		Source ProductSource `json:"_source"`
	}
	ProductSource struct {
		Title string  `json:"title"`
		Brand string  `json:"brand"`
		Price float64 `json:"price"`
		Stock int64   `json:"stock"`
	}
	FilterQuery struct {
		Term product.Filter `json:"term"`
	}
)

func New(es *es7.Client) product.Repository {
	return &Repository{es: es}
}

func (r *Repository) List(ctx context.Context, query string, filter product.Filter, sort *product.Sort, page *product.Page) (products []product.Product, resErr error) {
	products = []product.Product{}

	var q listQuery
	if page != nil {
		from := page.Size * (page.Num - 1)
		q.From = &from
		q.Size = &page.Size
	}
	if sort != nil {
		direction := "asc"
		if sort.Desc {
			direction = "desc"
		}
		sortBy := map[string]string{sort.Field: direction}
		q.Sort = []map[string]string{sortBy}
	}
	titleQuery := map[string]interface{}{
		"match": map[string]interface{}{"title": query},
	}

	if filter != nil {
		q.Query = map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   titleQuery,
				"filter": &FilterQuery{Term: filter},
			}}
	} else {
		q.Query = titleQuery
	}

	qJSON, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}

	resp, err := r.es.Search(
		r.es.Search.WithContext(context.Background()),
		r.es.Search.WithIndex("product"),
		r.es.Search.WithBody(bytes.NewReader(qJSON)),
		r.es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			resErr = err
		}
	}()

	if resp.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("error parsing the response body: %s", err)
		} else {
			return nil, fmt.Errorf("error from DB: %v", e["error"].(map[string]interface{})["reason"])
		}
	}

	var res responseModel
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("error parsing the product response body: %s", err)
	}

	for _, respModel := range res.Hits.Hits {
		products = append(products, respModel.Source.ToModel())
	}

	return products, nil
}

func (resp ProductSource) ToModel() product.Product {
	return product.Product{
		Brand: resp.Brand,
		Price: resp.Price,
		Stock: resp.Stock,
		Title: resp.Title,
	}
}

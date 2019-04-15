package elsearch

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"

	"github.com/kainobor/estest/config"
)

func New(conf config.Elastic) (*elasticsearch.Client, error) {
	esConf := elasticsearch.Config{Addresses: []string{"http://" + conf.IP + ":" + conf.Port}}
	es, err := elasticsearch.NewClient(esConf)
	if err != nil {
		return nil, fmt.Errorf("failed to create client for ES: %v", err)
	}

	return es, nil
}

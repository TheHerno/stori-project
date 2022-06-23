package spy

import (
	"github.com/elastic/go-elasticsearch/v7"
)

/*
Elasticsearch is a struct to inyect a spy on elasticsearch
*/
type Elasticsearch struct {
	Calls int //Number of calls
}

/*
NewElasticsearch is a constructor for Elasticsearch struct
*/
func NewElasticsearch() *Elasticsearch {
	return new(Elasticsearch)
}

/*
NewClient the signature of the function that's going to be mocked
*/
type NewClient func(cfg elasticsearch.Config) (*elasticsearch.Client, error)

/*
MockNewClient receives the mocked results, appends the arguments on each call
finally returns the mocked function
*/
func (spy *Elasticsearch) MockNewClient(firstResult *elasticsearch.Client, secondResult error) NewClient {
	return func(cfg elasticsearch.Config) (*elasticsearch.Client, error) {
		spy.Calls++
		return firstResult, secondResult
	}
}

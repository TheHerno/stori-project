package spy

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-extras/elogrus.v7"
)

/*
Elogrus is a struct to inyect a spy on elogrus
*/
type Elogrus struct {
	Calls int //Number of calls
}

/*
NewElogrus is a constructor for Elogrus struct
*/
func NewElogrus() *Elogrus {
	return new(Elogrus)
}

/*
NewAsyncElasticHookWithFunc the signature of the function that's going to be mocked
*/
type NewAsyncElasticHookWithFunc func(client *elasticsearch.Client, host string, level logrus.Level, indexFunc elogrus.IndexNameFunc) (*elogrus.ElasticHook, error)

/*
MockNewAsyncElasticHookWithFunc receives the mocked results, appends the arguments on each call
finally returns the mocked function
*/
func (spy *Elogrus) MockNewAsyncElasticHookWithFunc(firstResult *elogrus.ElasticHook, secondResult error) NewAsyncElasticHookWithFunc {
	return func(client *elasticsearch.Client, host string, level logrus.Level, indexFunc elogrus.IndexNameFunc) (*elogrus.ElasticHook, error) {
		spy.Calls++
		return firstResult, secondResult
	}
}

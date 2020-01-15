package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// servicePhonemizer is a phonemizer which uses a separate service.
//
// The phonemizers require high memory usage at least at least 128MB.
// Google App Engine F1 provides only 128MB of memory. So we need F2
// instead but it consumes twice instance hours. In the meanwhile, Heroku
// provides 512MB for the free plan.
//
// The phonemize has been separated on Heroku to serve api.hangulize.org for
// always free.
//
type servicePhonemizer struct {
	ctx context.Context
	id  string
}

func (p *servicePhonemizer) ID() string {
	return p.id
}

func (p *servicePhonemizer) URL(word string) string {
	var (
		baseURL     = "http://phonemize.herokuapp.com"
		sPhonemizer = url.PathEscape(p.id)
		sWord       = url.PathEscape(word)
	)

	// Override the phonemize base URL.
	switch p.ctx.Value("phonemizeURL").(type) {
	case string:
		baseURL = p.ctx.Value("phonemizeURL").(string)
	}

	return fmt.Sprintf("%s/%s/%s", baseURL, sPhonemizer, sWord)
}

func (p *servicePhonemizer) Phonemize(word string) string {
	phonemized, statusCode := p.PhonemizeStatusCode(word)
	if statusCode != 200 {
		log.Panicf("failed to phonemize %s/%s [%d]", p.id, word, statusCode)
	}
	return phonemized
}

func (p *servicePhonemizer) PhonemizeStatusCode(word string) (string, int) {
	url := p.URL(word)
	res, err := http.Get(url)

	if err != nil {
		return "", -1
	}

	var buf bytes.Buffer
	buf.ReadFrom(res.Body)

	return buf.String(), res.StatusCode
}

package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/url"

	"google.golang.org/appengine/urlfetch"
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
		hostname    = "phonemize.herokuapp.com"
		sPhonemizer = url.PathEscape(p.id)
		sWord       = url.PathEscape(word)
	)
	return fmt.Sprintf("http://%s/%s/%s", hostname, sPhonemizer, sWord)
}

func (p *servicePhonemizer) Phonemize(word string) string {
	c := urlfetch.Client(p.ctx)

	url := p.URL(word)
	res, err := c.Get(url)

	if err != nil {
		log.Panicf("failed to phonemize %s/%s: %s", p.id, word, err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(res.Body)
	return buf.String()
}

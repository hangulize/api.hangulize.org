package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/url"

	"google.golang.org/appengine/urlfetch"
)

type servicePhonemizer struct {
	ctx context.Context
	id  string
}

func (p *servicePhonemizer) ID() string {
	return p.id
}

func (p *servicePhonemizer) Phonemize(word string) string {
	// The phonemizers require high memory usage at least at least 128MB.
	// Google App Engine provides only 128MB of memory for free but Heroku
	// provides free 512MB. To always serve api.hangulize.org in free, the
	// phonemize service has been separated on Heroku.
	var (
		hostname    = "phonemize.herokuapp.com"
		sPhonemizer = url.PathEscape(p.id)
		sWord       = url.PathEscape(word)
	)
	url := fmt.Sprintf("http://%s/%s/%s", hostname, sPhonemizer, sWord)

	// Fetch the phonograms from a separate service.
	c := urlfetch.Client(p.ctx)
	res, err := c.Get(url)
	if err != nil {
		log.Panicf("failed to phonemize %s/%s: %s", p.id, word, err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(res.Body)
	return buf.String()
}

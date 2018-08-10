package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/url"

	"google.golang.org/appengine"
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
	hostname, err := appengine.ModuleHostname(p.ctx, "phonemize", "", "")
	if err != nil {
		log.Panicf("phonemize service hostname not found")
	}

	url := fmt.Sprintf(
		"http://%s/v2/phonemized/%s/%s",
		hostname,
		url.PathEscape(p.id),
		url.PathEscape(word),
	)

	// Fetch the phonograms from a separate service.
	c := urlfetch.Client(p.ctx)
	res, err := c.Get(url)
	if err != nil {
		log.Panicf("failed to get phonograms of %s/%s: %s", p.id, word, err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(res.Body)
	return buf.String()
}

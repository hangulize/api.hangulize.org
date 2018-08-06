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

type servicePronouncer struct {
	ctx context.Context
	id  string
}

func (p *servicePronouncer) ID() string {
	return p.id
}

func (p *servicePronouncer) Pronounce(word string) string {
	hostname, err := appengine.ModuleHostname(p.ctx, p.id, "", "")
	if err != nil {
		log.Panicf("service hostname for %s not found", p.id)
	}

	url := fmt.Sprintf(
		"http://%s/v2/pronounced/%s/%s",
		hostname,
		url.PathEscape(p.id),
		url.PathEscape(word),
	)

	// Fetch the pronunciation from a separate service.
	c := urlfetch.Client(p.ctx)
	res, err := c.Get(url)
	if err != nil {
		log.Panicf("failed to get pronunciation of %s/%s: %s", p.id, word, err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(res.Body)
	return buf.String()
}

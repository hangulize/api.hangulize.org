package main

import (
	"github.com/hangulize/hangulize/pronounce/furigana"
	"github.com/hangulize/hangulize/pronounce/pinyin"
)

// Copied from hangulize/pronouncer.go to not import hangulize.
type pronouncer interface {
	ID() string
	Pronounce(string) string
}

var pronouncers = map[string]pronouncer{
	furigana.P.ID(): &furigana.P,
	pinyin.P.ID():   &pinyin.P,
}

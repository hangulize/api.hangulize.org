package main

import (
	"github.com/hangulize/hangulize/phonemize/furigana"
	"github.com/hangulize/hangulize/phonemize/pinyin"
)

// Copied from hangulize/hangulizePhonemizer.go to not import hangulize.
type hangulizePhonemizer interface {
	ID() string
	Phonemize(string) string
}

var phonemizers = map[string]hangulizePhonemizer{
	furigana.P.ID(): &furigana.P,
	pinyin.P.ID():   &pinyin.P,
}

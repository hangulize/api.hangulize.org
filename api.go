package main

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/hangulize/hangulize"
)

// Register adds routing rules onto the given router group.
func Register(r gin.IRouter) {
	r.GET("/version", Version)
	r.GET("/specs", Specs)
	r.GET("/specs/:path", SpecHGL)
	r.GET("/hangulized/:lang/*word", Hangulized)
	r.GET("/phonemized/:phonemizer/*word", Phonemized)
}

// Version returns the version of the "hangulize" package.
//
//  Route:  GET /version
//  Accept: text/plain (default), application/json
//
func Version(c *gin.Context) {
	switch c.NegotiateFormat(gin.MIMEPlain, gin.MIMEJSON) {

	case gin.MIMEJSON:
		c.JSON(http.StatusOK, gin.H{
			"version": hangulize.Version,
		})

	case gin.MIMEPlain:
		fallthrough

	default:
		c.String(http.StatusOK, hangulize.Version)
	}
}

// Specs returns the list of supported specs.
//
//  Route:  GET /specs
//  Accept: text/plain (default), application/json
//
func Specs(c *gin.Context) {
	switch c.NegotiateFormat(gin.MIMEPlain, gin.MIMEJSON) {

	case gin.MIMEJSON:
		langs := hangulize.ListLangs()
		specs := make([]*gin.H, len(langs))

		for i, lang := range langs {
			spec, _ := hangulize.LoadSpec(lang)
			specs[i] = packSpec(spec)
		}

		c.JSON(http.StatusOK, gin.H{
			"specs": specs,
		})

	case gin.MIMEPlain:
		fallthrough

	default:
		c.String(http.StatusOK, strings.Join(hangulize.ListLangs(), "\n"))
	}
}

func packSpec(s *hangulize.Spec) *gin.H {
	test := make([]gin.H, len(s.Test))

	for i, example := range s.Test {
		test[i] = gin.H{
			"word":        example[0],
			"transcribed": example[1],
		}
	}

	return &gin.H{
		"lang": gin.H{
			"id":         s.Lang.ID,
			"codes":      s.Lang.Codes,
			"english":    s.Lang.English,
			"korean":     s.Lang.Korean,
			"script":     s.Lang.Script,
			"phonemizer": s.Lang.Phonemizer,
		},

		"config": gin.H{
			"authors": s.Config.Authors,
			"stage":   s.Config.Stage,
		},

		"test": test,
	}
}

// SpecHGL serves the HGL source of the spec.
//
//  Route:  GET /specs/{lang}.hgl
//  Accept: text/vnd.hgl
//
func SpecHGL(c *gin.Context) {
	// Should look like "ita.hgl".
	path := c.Param("path")

	if !strings.HasSuffix(path, ".hgl") {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	lang := strings.TrimSuffix(path, ".hgl")
	spec, ok := hangulize.LoadSpec(lang)

	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Header("Content-Type", "text/vnd.hgl")
	c.String(http.StatusOK, spec.Source)
}

// paramWord picks the word argument.
//
// The parameter should be declared as *word.
//
func paramWord(c *gin.Context) string {
	word := c.Param("word")

	// Remove the initial slash.
	word = word[1:]

	unescaped, err := url.QueryUnescape(word)
	if err != nil {
		word = unescaped
	}

	return word
}

// Hangulized transcribes a non-Korean word into Hangul.
//
//  Route:  GET /hangulized/{lang}/{word}
//  Accept: text/plain (default), application/json
//
func Hangulized(c *gin.Context) {
	lang := c.Param("lang")
	word := paramWord(c)

	spec, ok := hangulize.LoadSpec(lang)
	defer hangulize.UnloadSpec(lang)
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	h := hangulize.NewHangulizer(spec)

	phonemizer := spec.Lang.Phonemizer
	if phonemizer != "" {
		ctx := appengine.NewContext(c.Request)
		h.UsePhonemizer(&servicePhonemizer{ctx, phonemizer})
	}

	transcribed := h.Hangulize(word)

	switch c.NegotiateFormat(gin.MIMEPlain, gin.MIMEJSON) {

	case gin.MIMEJSON:
		c.JSON(http.StatusOK, gin.H{
			"lang":        lang,
			"word":        word,
			"transcribed": transcribed,
		})

	case gin.MIMEPlain:
		fallthrough

	default:
		c.String(http.StatusOK, transcribed)
	}
}

// Phonemized guesses phonograms from a spelling.
//
//  Route:  GET /phonemized/{phonemizer}/{word}
//
func Phonemized(c *gin.Context) {
	pID := c.Param("phonemizer")
	word := paramWord(c)

	ctx := appengine.NewContext(c.Request)
	p := servicePhonemizer{ctx, pID}

	phonemized, statusCode := p.PhonemizeStatusCode(word)

	if statusCode != 200 {
		c.AbortWithStatus(statusCode)
		return
	}

	switch c.NegotiateFormat(gin.MIMEPlain, gin.MIMEJSON) {

	case gin.MIMEJSON:
		c.JSON(http.StatusOK, gin.H{
			"phonemizer": pID,
			"word":       word,
			"phonemized": phonemized,
		})

	case gin.MIMEPlain:
		fallthrough

	default:
		c.String(http.StatusOK, phonemized)
	}
}

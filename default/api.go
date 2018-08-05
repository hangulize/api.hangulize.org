package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/hangulize/hangulize"
)

// Register adds routing rules onto the given router group.
func Register(r gin.IRouter) {
	r.GET("/version", Version)
	r.GET("/specs", Specs)
	r.GET("/specs/:path", SpecHGL)
	r.GET("/hangulized/:lang/:word", Hangulized)
	r.GET("/pronounced/:_/:_", Pronounced)
}

// Version returns the version of the "hangulize" package.
//
//  Route:  GET /version
//  Accept: text/plain (default), application/json
//
func Version(c *gin.Context) {
	switch c.NegotiateFormat(gin.MIMEJSON) {

	case gin.MIMEJSON:
		c.JSON(http.StatusOK, gin.H{
			"version": hangulize.Version,
		})

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
	switch c.NegotiateFormat(gin.MIMEJSON) {

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
			"id":        s.Lang.ID,
			"codes":     s.Lang.Codes,
			"english":   s.Lang.English,
			"korean":    s.Lang.Korean,
			"script":    s.Lang.Script,
			"pronounce": s.Lang.Pronounce,
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

// Hangulized transcribes a non-Korean word into Hangul.
//
//  Route:  GET /hangulized/{lang}/{word}
//  Accept: text/plain (default), application/json
//
func Hangulized(c *gin.Context) {
	lang := c.Param("lang")
	word := c.Param("word")

	transcribed := hangulize.Hangulize(lang, word)

	switch c.NegotiateFormat(gin.MIMEJSON) {

	case gin.MIMEJSON:
		c.JSON(http.StatusOK, gin.H{
			"lang":        lang,
			"word":        word,
			"transcribed": transcribed,
		})

	default:
		c.String(http.StatusOK, transcribed)
	}
}

// Pronounced guesses a pronunciation from a spelling. But each pronouncer
// should be implemented in a separate service for cost efficiency. This
// handler always serves the 404 error.
//
//  Route:  GET /pronounced/{pronouncer}/{word}
//
func Pronounced(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotFound)
}

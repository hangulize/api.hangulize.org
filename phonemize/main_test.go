package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func init() {
	router = gin.New()
	register(router)
}

func GET(path string, accept string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	req.Header.Set("Accept", accept)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	return w
}

func TestPhonemizedText(t *testing.T) {
	w := GET("/phonemized/furigana/自由ヶ丘", "text/plain")

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/plain")

	assert.Equal(t, "ジユウガオカ", w.Body.String())
}

func TestPhonemizedJSON(t *testing.T) {
	w := GET("/phonemized/furigana/自由ヶ丘", "application/json")

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var res map[string]string
	json.Unmarshal(w.Body.Bytes(), &res)

	expected := map[string]string{
		"phonemizer": "furigana",
		"word":       "自由ヶ丘",
		"phonemized": "ジユウガオカ",
	}
	assert.Equal(t, expected, res)
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"google.golang.org/appengine/aetest"

	"github.com/hangulize/hangulize"
)

var router *gin.Engine
var inst aetest.Instance

func init() {
	router = gin.New()
	Register(router)

	_inst, err := aetest.NewInstance(nil)
	if err != nil {
		panic(err)
	}
	inst = _inst
}

// GET sends a GET request.
func GET(path, accept string) *httptest.ResponseRecorder {
	return GETWithValue(path, accept, nil, nil)
}

// GETWithValue sends a GET request with a context value.
func GETWithValue(
	path, accept string,
	key, val interface{},
) *httptest.ResponseRecorder {
	req, _ := inst.NewRequest("GET", path, nil)
	req.Header.Set("Accept", accept)

	if key != nil {
		ctx := context.WithValue(req.Context(), key, val)
		req = req.WithContext(ctx)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestVersionText(t *testing.T) {
	w := GET("/version", "text/plain")

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/plain")

	assert.Equal(t, hangulize.Version, w.Body.String())
}

func TestVersionJSON(t *testing.T) {
	w := GET("/version", "application/json")

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var res map[string]string
	json.Unmarshal(w.Body.Bytes(), &res)

	expected := map[string]string{
		"version": hangulize.Version,
	}
	assert.Equal(t, expected, res)
}

func TestSpecsText(t *testing.T) {
	w := GET("/specs", "text/plain")

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/plain")

	firstLang := hangulize.ListLangs()[0]
	assert.True(t, strings.HasPrefix(w.Body.String(), firstLang+"\n"))
}

func TestSpecsJSON(t *testing.T) {
	w := GET("/specs", "application/json")

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var res map[string][]struct {
		Lang struct {
			ID      string
			Codes   []string
			English string
			Korean  string
			Script  string
		}

		Config struct {
			Authors []string
			Stage   string
		}

		Test []struct {
			Word        string
			Transcribed string
		}
	}
	json.Unmarshal(w.Body.Bytes(), &res)
	firstRes := res["specs"][0]

	firstLang := hangulize.ListLangs()[0]
	firstSpec, _ := hangulize.LoadSpec(firstLang)

	assert.Equal(t, firstLang, firstRes.Lang.ID)
	assert.Equal(t, firstSpec.Lang.English, firstRes.Lang.English)
	assert.Equal(t, firstSpec.Config.Stage, firstRes.Config.Stage)
	assert.Equal(t, firstSpec.Test[0][0], firstRes.Test[0].Word)
	assert.Equal(t, firstSpec.Test[0][1], firstRes.Test[0].Transcribed)
}

func TestHangulizedText(t *testing.T) {
	w := GET("/hangulized/ita/gloria", "text/plain")

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/plain")

	assert.Equal(t, "글로리아", w.Body.String())
}

func TestHangulizedJSON(t *testing.T) {
	w := GET("/hangulized/ita/gloria", "application/json")

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	var res map[string]string
	json.Unmarshal(w.Body.Bytes(), &res)

	expected := map[string]string{
		"lang":        "ita",
		"word":        "gloria",
		"transcribed": "글로리아",
	}
	assert.Equal(t, expected, res)
}

func TestHangulizedIntegratedWithPhonemize(t *testing.T) {
	w := GET("/hangulized/jpn/自由ヶ丘", "text/plain")

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/plain")

	assert.Equal(t, "지유가오카", w.Body.String())
}

func TestHangulizedSlash(t *testing.T) {
	r := GET("/hangulized/ita/gloria%2Fgloria", "text/plain")
	assert.Equal(t, "글로리아/글로리아", r.Body.String())
}

func newPhonemizeServer() *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/furigana/") {
				fmt.Fprintln(w, "ジユーガオカ")
			} else {
				w.WriteHeader(404)
			}
		}),
	)
}

func TestPhonemized(t *testing.T) {
	ps := newPhonemizeServer()
	defer ps.Close()

	w := GETWithValue(
		"/phonemized/furigana/自由ヶ丘", "text/plain",
		"phonemizeURL", ps.URL,
	)
	assert.Equal(t, 200, w.Code)
	assert.NotEqual(t, "", w.Body.String())
}

func TestPhonemized404(t *testing.T) {
	ps := newPhonemizeServer()
	defer ps.Close()

	w := GETWithValue(
		"/phonemized/unknown/nwonknu", "text/plain",
		"phonemizeURL", ps.URL,
	)
	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "", w.Body.String())
}

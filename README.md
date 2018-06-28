# api.hangulize.org

Web API for Hangulize, powered by Google App Engine

## v2

`v2` is the re-designed Hangulize Web API. It is much sane than `v1`.

### GET /v2/version

Returns the version of the underlying Hangulize.

| Content-Type       | Result Example  |
| ------------------ | --------------- |
| `text/plain`       | `0.1.0`         |
| `application/json` | `{"version": "0.1.0"}` |

### GET /v2/hangulized/`:lang`/`:word`

Transcribes a non-Korean word into Hangul.
The result of `GET /v2/hangulized/ita/gloria` would be:

| Content-Type       | Result Example |
| ------------------ | -------------- |
| `text/plain`       | `글로리아`      |
| `application/json` | `{"lang": "ita", "word": "gloria", "transcribed": "글로리아"}` |

### GET /v2/specs

Provides the list of language-specific transcription specs.

| Content-Type       | Result Example |
| ------------------ | -------------- |
| `text/plain`       | `ita jpn deu...` |
| `application/json` | `{"specs": [{"lang": {"id": "ita", "english": "Italian", ...}, ...]}` |

### GET /v2/specs/`:lang`.hgl

The source code of the spec.
The only result format is `text/vnd.hgl`.

```hgl
lang:
    id      = "ita"
    codes   = "it", "ita"
    english = "Italian"
    korean  = "이탈리아어"
    script  = "latin"
```

## v1

`v1` re-implements [the original Hangulize Web API](https://pythonhosted.org/hangulize/webapi.html)
except some unnecessary specs.

## Development

To run locally, use Cloud SDK:

```console
$ dev_appserver.py app.yaml
```

Or just run the `main` package:

```console
$ go run *.go
```

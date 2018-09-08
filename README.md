# api.hangulize.org

Web API for Hangulize, powered by Google App Engine

## API

You can find the actual API specification in [openapi.yaml](openapi.yaml) under
OpenAPI Specification 3.0.1.

There was a former Hangulize Web API written in Python. That is considered as
`v1`. The code still can be found at https://github.com/sublee/hangulize-web.
The `v1` API had complex content negotiation rules and not deterministic
behaviors. `v2` is re-designed against `v1` to provide more simple and
predictable usage.

### GET /v2/version

Returns the version of the underlying Hangulize.

| Content-Type       | Result Example  |
| ------------------ | --------------- |
| `text/plain`       | `0.1.0`         |
| `application/json` | `{"version": "0.1.0"}` |

### GET /v2/specs

Provides the list of language-specific transcription specs.

| Content-Type       | Result Example |
| ------------------ | -------------- |
| `text/plain`       | `ita jpn deu...` |
| `application/json` | `{"specs": [{"lang": {"id": "ita", "english": "Italian", ...}, ...]}` |

### GET /v2/specs/`{lang}`.hgl

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

### GET /v2/hangulized/`{lang}`/`{word}`

Transcribes a non-Korean word into Hangul.
The result of `GET /v2/hangulized/ita/cappuccino` would be:

| Content-Type       | Result Example |
| ------------------ | -------------- |
| `text/plain`       | `글로리아`      |
| `application/json` | `{"lang": "ita", "word": "cappuccino", "transcribed": "카푸치노"}` |

### GET /v2/phonemized/`{phonemizer}`/`{word}`

Guesses phonograms from a spelling.
The result of `GET /v2/phonemized/furigana/東京` would be:

| Content-Type       | Result Example |
| ------------------ | -------------- |
| `text/plain`       | `トーキョー`      |
| `application/json` | `{"phonemizer": "furigana", "word": "東京", "phonemized": "トーキョー"}` |

## Development

To run locally, use Cloud SDK:

```console
$ dev_appserver.py app.yaml
```

Or just build the `main` package:

```console
$ go build && ./api.hangulize.org
```

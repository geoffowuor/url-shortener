# URL Shortener

A simple URL shortening API built with Go, Fiber v3, GORM, and SQLite.

## Features

- Create short URLs
- Automatic short-code generation
- URL redirects
- Click tracking
- URL statistics
- SQLite persistence

## API

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/shorten` | Create a short URL |
| GET | `/:shortCode` | Redirect to original URL |
| GET | `/api/stats/:shortCode` | Get URL statistics |

### Create URL

```json
{
  "original_url": "https://github.com/golang/go"
}
```

Example response:

```json
{
  "id": 1,
  "original_url": "https://github.com/golang/go",
  "short_code": "BvIf4n",
  "clicks": 0
}
```

## Stack

- Go
- Fiber v3
- GORM
- SQLite

## Run

```bash
go mod download
go run .
```

Server: `http://localhost:3000`

## Structure

```
url-shortener/
├── internals/
│   ├── handlers/
│   ├── models/
│   └── utils/
├── main.go
├── go.mod
```
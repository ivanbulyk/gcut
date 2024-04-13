# gCUT 

## URL shortening service 

## Prerequisites

- Go 1.22 or higher
- Redis server
- Docker (optional for containerized deployment)

## Getting Started

### Clone the Repository

```bash
git clone https://github.com/ivanbulyk/gcut.git
cd gcut
```

### Set Up Environment Variables

Create a `.env` file at the root of the project and add the following content:

```txt
REDIS_ADDRESS=redis:6379
REDIS_PASS=
RATE_LIMIT=10
RATE_LIMIT_RESET=30m
HOST=0.0.0.0
PORT=8082



DOMAIN=localhost:8082
```

Replace `<HOST>,<PORT>` with actual values when needed. The same goes for other environment variables

### Running with Docker

Build and start the container:

```bash
docker compose up -d --build
```

Track the logs:

```bash
docker compose logs -f
```

Hit `http://127.0.0.1:8082` to health check.

Make POST request to `http://127.0.0.1:8082/api/v1/` with body containing following payload:
{"url":"your_long_url_you_want_to_shorten"} to get shorten version.

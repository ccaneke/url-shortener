# URL Shortener: A Go implementation of a URL shortener HTTP service
This URL shortener HTTP service shows short, easy to type URLs to end users.
### What it does?
* Provides an HTTP API to:
    * Shorten a URL
    * Redirect to the long URL from the shortened URL
### Shortened URL characteristics:
* The ID of the shortened URL is unique (across past and concurrent requests)
* The ID of the shortened URL is short (max. 10 characters long)
* The long/shortened URL mapping is persisted and isn't lost after a backend service restart. For example when this service is run under the domain https://short.com this URL https://en.wikipedia.org/wiki/Main_Page/123456 should look like this https://short.com/c9eo1 after being shortened. When calling https://short.com/c9eo1 a user should be redirected to https://en.wikipedia.org/wiki/Main_Page/123456.

## Setup and how to run

After cloning the repository, run the project by executing the following command in the same directory as the `docker-compose.yml` file, making sure to pass in the redis database password as an environment variable:

```
sudo REDIS_DATABASE_PASSWORD=changetheworld docker compose up
```

Alternatively, you can run the project with `go run .` in the `cmd/httpapi/server` directory or install it by running `go install` in that same directory and then use:

```
url-shortner-poc
```

**Note that to run the project with docker compose you must install [docker](https://docs.docker.com/engine/install/) and [docker compose](https://docs.docker.com/compose/install/). Other wise, if you're going to run the project on your local machine you must have a `redis-server` running in the background and set `appendonly yes` in the `redis.conf` file. See [instructions](https://redis.io/docs/getting-started/installation/) for installing redis. On Linux you can find the `redis.conf` file in `/etc/redis/redis.conf` and on mac `/usr/local/etc/redis.conf`.

## Routes
### `localhost:8080/shorten`:
POST a JSON object like the one below, and receive the short URL in the response:
```JSON
{
    "LongURL": "https://en.wikipedia.org/wiki/Main_Page/123456"
}
```

Result:
```JSON
{
    "OriginalURL": "https://en.wikipedia.org/wiki/Main_Page/123456",
    "ShortURL": "https://localhost/04ec8683"
}
```
### `localhost:8080/{YOUR-SHORTENED-URL}`
Send a GET request and you will be redirected to your long URL
 
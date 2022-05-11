# URL Shortener

## Introduction

This is a URL shortening API written in Go. It supports the following high level requirements:

* Creating short urls with either a custom slug or random one
* Deleting short urls
* Getting simple statistics about short urls
* Accessing short urls

## Running

The easiest way to run the application is to use `docker-compose`. You should run the following to get the application started:

```
docker compose build
docker compose up
```

or you can use the included `makefile`:

```
make docker-up
```

This will start the webserver on `localhost:8080` and you can begin to issue requests to the API.

## Routes

The application exposes the following routes:

| HTTP Verb     | Route               | Description |
| ------------- | --------------------| ----------- |
| `GET`         | `/:slug`            | Access a short URL. Clients are redirected to the long url associated with the given slug
| `POST`        | `/api/v1/shorturls` | Create a new short URL. Clients can specify their own custom slug or let the system generate a random one.
| `DELETE`      | `/api/v1/shorturls/:slug` | Delete the short URL associated with the given slug
| `GET`         | `/api/v1/shorturls/:slug` | Get short URL information associated with the given slug
| `GET`         | `/api/v1/shorturls/:slug/clicks` | Get analytics data associated with the given slug

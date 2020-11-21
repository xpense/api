# xpense API

## Introduction

The **xpense API** allows users to keep track of their income and expenses. Users can create `transactions` that can be grouped in `wallets` and associated with different `parties`.

## Local development

### Prerequisites

Installation instructions for the packages below are for `macOS`.
For `Linux` and `Windows` users - please use the linked resources to download the packages for your operating system of choice.

- **golang**

  A working [go](https://golang.org/dl/) installation. The latest stable version (at the time of writing `1.15.x`) is recommended.

  With Homebrew:

  ```sh
  brew install go
  ```

- **air**

  For live reloading during development: [cosmtrek/air](https://github.com/cosmtrek/air)

  ```sh
  go get -u github.com/cosmtrek/air
  ```

- **mockery**

  Needed for (re)generating mock files for testing: [vektra/mockery](https://github.com/vektra/mockery)

  With Homebrew:

  ```sh
  brew install mockery
  ```

- **Postgres**

  You need a running [PostgreSQL](https://www.postgresql.org/download/) server to use the API.

  `Docker` users can run the provided `docker-compose.yml` file.

  ```sh
  # create a volume first to persist data
  docker volume create xpense_postgresql

  # start the docker container
  docker-compose up
  ```

- **Postman (optional)**

  There is an exported Postman Collection in `/docs` that you can use to make requests locally.

  1. Download [Postman](https://www.postman.com/downloads/).
  2. Import the collection from `/docs`
  3. Use the `/auth/signup` method to register
  4. Use the `/auth/login` method to login with your newly created account (an authentication token will be saved automatically for you and you will be able to make subsequent requests to all routes without having to set the `Authorization` header)

### Setting up environment variables

Copy the example `env` file and set the required variables

```sh
cp .env.example.yml .env.yml
```

## Running the dev server

If you installed `air` for live reloading:

```sh
# inside the top-level directory
air
```

Standard `go` way (no live-reloading):

```sh
go run main.go
```

## Running the test suite

To run all the tests:

```sh
go test ./...
```

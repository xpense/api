# xpense API

- [xpense API](#xpense-api)
  - [Introduction](#introduction)
  - [Local development](#local-development)
    - [Prerequisites](#prerequisites)
    - [Setting up environment variables](#setting-up-environment-variables)
    - [Running the dev server](#running-the-dev-server)
    - [Running the test suite](#running-the-test-suite)
    - [Building the dockerized application](#building-the-dockerized-application)
  - [Backend Architecture](#backend-architecture)
  - [Database](#database)
  - [Documentation](#documentation)
    - [Authentication](#authentication)
      - [Sign Up](#sign-up)
      - [Login](#login)
    - [Account](#account)
      - [Get Account information](#get-account-information)
      - [Update Account information](#update-account-information)
      - [Delete Account](#delete-account)
    - [Wallets](#wallets)
      - [Create Wallet](#create-wallet)
      - [Get Wallet](#get-wallet)
      - [Update Wallet](#update-wallet)
      - [Delete Wallet](#delete-wallet)
      - [List Wallets](#list-wallets)
      - [List Transactions by Wallet](#list-transactions-by-wallet)
    - [Parties](#parties)
      - [Create Party](#create-party)
      - [Get Party](#get-party)
      - [Update Party](#update-party)
      - [Delete Party](#delete-party)
      - [List Parties](#list-parties)
      - [List Transactions by Party](#list-transactions-by-party)
    - [Transactions](#transactions)
      - [Create Transaction](#create-transaction)
      - [Get Transaction](#get-transaction)
      - [Update Transaction](#update-transaction)
      - [Delete Transaction](#delete-transaction)
      - [List all Transactions](#list-all-transactions)
  - [Contributors](#contributors)

## Introduction

The xpense API allows users to keep track of their income and expenses. Users can create transactions that can be grouped in wallets and associated with different parties.

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

  `Docker` users can run the provided `docker-compose.dev.yml` file.

  ```sh
  # start the dev docker container
  docker-compose -f docker-compose.dev.yml up
  ```

- **Postman (optional)**

  There is an exported Postman Collection in `/docs` that you can use to make requests locally.

  1. Download [Postman](https://www.postman.com/downloads/).
  2. Import the collection from `/docs`
  3. Use the `/api/v1/auth/signup` method to register
  4. Use the `/api/v1/auth/login` method to login with your newly created account (an authentication token will be saved automatically for you and you will be able to make subsequent requests to all routes without having to set the `Authorization` header)

### Setting up environment variables

Copy the example `env` file and set the required variables

```sh
cp .env.example .env
```

### Running the dev server

Run the following command inside the top-level directory:

- If you installed `air` for live reloading:

  ```sh
  air
  ```

- Standard `go` way (no live-reloading):

  ```sh
  go run main.go
  ```

### Running the test suite

To run all the tests:

```sh
go test ./...
```

### Building the dockerized application

To build the application together with a running database you can use the provided `docker-compose.yml` file.

```sh
docker-compose up --build
```

## Backend Architecture

![Architecture Diagram](docs/XpenseBackend.png)

The live version is running on Digital Ocean Ubuntu 20.4 Droplet. Nginx serves as a reverse proxy within Docker and redirects HTTP traffic to HTTPS and adds SSL credentals, protocols, and security headers. It redirects all incoming requests to the containerized golang application. Internally the application also has access to the Postgres DB also running in a Docker container. All containers are running within the same network.

## Database

![ER-model](docs/ER.png)

## Documentation

The backend is versioned. The current version is `v1`.

All endpoints may return `500 Internal Server Error`, when something unexpected happens on the server side.

### Authentication

The API uses the [JWT standard](https://jwt.io/) to authenticate users and protect resources and routes

#### Sign Up

Endpoint:

```text
POST /api/v1/auth/signup
```

Request payload:

```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "email@example.com",
  "password": "123Password!{}"
}
```

`password` has to contain at least 1 upper, lower, numerical and special char and have a minimum length of 8 chars

Responses:

- `201 Created`

  Account was succesfully created

- `400 Bad Request`

  Somethinig went wrong when processing the request. Either empty request body, malformed request body, missing first or last name, or invalid email/password.

- `409 Conflict`

  If you're trying to create a new account with an already registered email.

#### Login

Endpoint:

```text
POST /api/v1/auth/login
```

Request payload:

```json
{
  "email": "email@example.com",
  "password": "123Password!{}"
}
```

Responses:

- `200 OK`

  Account was logged in succesfully

  Example:

  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
  }
  ```

- `400 Bad Request`

  Somethinig went wrong when processing the request. Either empty request body, malformed request body, missing email or password field, or the provided password is wrong.

- `404 Not Found`

  Account with the specified email does not exist.

### Account

All routes are protected and require the following header with a valid authentication token (can be obtained from [Login](#login)):

```text
Authorization: Bearer <token>
```

#### Get Account information

Endpoint:

```text
GET /api/v1/account
```

Responses:

- `200 OK`

  Account was retrieved succesfully

  Example:

  ```json
  {
    "id": "1",
    "created_at": "2020-11-20T15:05:36.248855+01:00",
    "updated_at": "2020-11-20T15:05:36.248855+01:00",
    "first_name": "Jane",
    "last_name": "Doe",
    "email": "jane@doe.com"
  }
  ```

- `400 Bad Request`

  Somethinig went wrong when processing the request. Either empty request body, malformed request body, missing email or password field, or the provided password is wrong.

- `401 Unauthorized`

  The provided token is not valid.

- `404 Not Found`

  Account with the ID belonging to the token does not exist (possibly deleted).

#### Update Account information

Endpoint:

```text
PATCH /api/v1/account
```

Request payload:

```json5
{
  "first_name": "Jenny",       // optional
  "last_name": "Doh",          // optional
  "email": "jenny@doh.com"     // optional
}
```

Responses:

- `200 OK`

  Acccount was updated successfully.

  Example:

  ```json
  {
    "id": "1",
    "created_at": "2020-11-20T15:05:36.248855+01:00",
    "updated_at": "2020-11-20T15:05:36.248855+01:00",
    "first_name": "Jenny",
    "last_name": "Doh",
    "email": "jenny@doh.com"
  }
  ```

- `400 Bad Request`

  Somethinig went wrong when processing the request. Either empty request body, malformed request body, or invalid email.

- `401 Unauthorized`

  The provided token is not valid.

- `404 Not Found`

  Account with the ID belonging to the token does not exist (possibly deleted).

#### Delete Account

Endpoint:

```text
DELETE /api/v1/account
```

Responses:

- `204 No Content`

  Acccount was deleted successfully.

- `401 Unauthorized`

  The provided token is not valid.

- `404 Not Found`

  Account with the ID belonging to the token does not exist (possibly previously deleted).

### Wallets

A wallet represents a group of transactions belonging to a user. One user can have multiple wallets (e.g. one for cash, one for the bank, one for work)

All routes are protected and require the following header with a valid authentication token (can be obtained from [Login](#login)):

```text
Authorization: Bearer <token>
```

#### Create Wallet

Endpoint:

```text
POST /api/v1/wallets
```

Request payload:

```json5
{
  "name": "cash",
  "description": "a wallet only for cash transactions" // optional
}
```

Responses:

- `201 Created`

  Wallet was created successfully.

  Example:

  ```json
  {
    "id": 2,
    "created_at": "2020-11-20T15:06:27.277849+01:00",
    "updated_at": "2020-11-20T15:06:27.277849+01:00",
    "name": "cash",
    "description": "a wallet only for cash transactions"
  }
  ```

- `400 Bad Request`

  Somethinig went wrong when processing the request. Either empty request body or malformed request body.

- `401 Unauthorized`

  The provided token is not valid.

- `409 Conflict`

  A wallet with the same name belonging to the same user already exists.

#### Get Wallet

Endpoint:

```text
GET /api/v1/wallets/:id
```

where `:id` is the ID of the wallet you want to retrieve

Responses:

- `200 OK`

  Wallet was retrieved successfully.

  Example:

  ```json
  {
    "id": 2,
    "created_at": "2020-11-20T15:06:27.277849+01:00",
    "updated_at": "2020-11-20T15:06:27.277849+01:00",
    "name": "cash",
    "description": "a wallet only for cash transactions"
  }
  ```

- `401 Unauthorized`

  The provided token is not valid.

- `403 Forbidden`

  The wallet with the specified ID does not belong to the current user.

- `404 Not Found`

  The wallet with the specified ID does not exist.

#### Update Wallet

Endpoint:

```text
PATCH /api/v1/wallets/:id
```

where `:id` is the ID of the wallet you want to update

Request payload:

```json5
{
  "name": "Cash",                   // optional
  "description": "my cash wallet"   // optional
}
```

Responses:

- `200 OK`

  Wallet was updated successfully.

  Example:

  ```json
  {
    "id": 2,
    "created_at": "2020-11-20T15:06:27.277849+01:00",
    "updated_at": "2020-11-20T19:20:14.277849+01:00",
    "name": "Cash",
    "description": "my cash wallet"
  }
  ```

- `400 Bad Request`

  Somethinig went wrong when processing the request. Either empty request body or malformed request body.

- `401 Unauthorized`

  The provided token is not valid.

- `403 Forbidden`

  The wallet with the specified ID does not belong to the current user.

- `404 Not Found`

  The wallet with the specified ID does not exist.

- `409 Conflict`

  A wallet with the same name belonging to the same user already exists.

#### Delete Wallet

Endpoint:

```text
DELETE /api/v1/wallets/:id
```

where `:id` is the ID of the wallet you want to delete

Responses:

- `204 No Content`

  Wallet was deleted successfully.

- `401 Unauthorized`

  The provided token is not valid.

- `403 Forbidden`

  The wallet with the specified ID does not belong to the current user.

- `404 Not Found`

  The wallet with the specified ID does not exist.

#### List Wallets

Lists all wallets which belong to the currently logged-in user.

Endpoint:

```text
GET /api/v1/wallets
```

Responses:

- `200 OK`

  Wallets were retrieved successfully.

  Example:

  ```json
  {
    "count": 2,
    "entries": [
      {
        "id": 2,
        "created_at": "2020-11-20T15:06:27.277849+01:00",
        "updated_at": "2020-11-20T15:06:27.277849+01:00",
        "name": "cash",
        "description": "a wallet for only cash transactions"
      },
      {
        "id": 5,
        "created_at": "2020-11-20T15:11:44.906804+01:00",
        "updated_at": "2020-11-20T15:12:19.46906+01:00",
        "name": "Sparkasse",
        "description": "a wallet for banking transactions"
      }
    ]
  }
  ```

- `401 Unauthorized`

  The provided token is not valid.

- `404 Not Found`

  The wallet with the specified ID does not exist.

#### List Transactions by Wallet

Lists all transactions which are associated with the wallet with the specified ID.

Endpoint:

```text
GET /api/v1/wallets/:id/transactions
```

Responses:

- `200 OK`

  Transactions were retrieved successfully.

  Example:

  ```json
  {
    "count": 1,
    "entries": [
      {
        "id": 4,
        "wallet_id": 2,
        "party_id": 2,
        "created_at": "2020-11-20T17:06:47.329695+01:00",
        "updated_at": "2020-11-20T17:06:47.329695+01:00",
        "timestamp": "2020-11-20T17:06:47.327376+01:00",
        "amount": "29.99",
        "description": "groceries"
      },
      {
        "id": 10,
        "wallet_id": 2,
        "party_id": 3,
        "created_at": "2020-11-20T17:06:47.329695+01:00",
        "updated_at": "2020-11-20T17:06:47.329695+01:00",
        "timestamp": "2020-11-20T17:06:47.327376+01:00",
        "amount": "6.99",
        "description": "cigarettes"
      }
    ]
  }
  ```

- `401 Unauthorized`

  The provided token is not valid.

- `403 Forbidden`

  The wallet with the specified ID does not belong to the current user.

### Parties

A party represents the sender or the recipient of a transaction created by the user. If the transaction is an expense, then the party represents the recipient, whereas if the transaction is an income, then the party represents the sender.

All routes are protected and require the following header with a valid authentication token (can be obtained from [Login](#login)):

```text
Authorization: Bearer <token>
```

#### Create Party

Endpoint:

```text
POST /api/v1/parties
```

Request payload:

```json5
{
  "name": "Amazon",
}
```

Responses:

- `201 Created`

  Party was created successfully.

  Example:

  ```json
  {
    "id": 3,
    "created_at": "2020-11-20T15:06:27.277849+01:00",
    "updated_at": "2020-11-20T15:06:27.277849+01:00",
    "name": "Amazon",
  }
  ```

- `400 Bad Request`

  Somethinig went wrong when processing the request. Either empty request body or malformed request body.

- `401 Unauthorized`

  The provided token is not valid.

- `409 Conflict`

  A party with the same name belonging to the same user already exists.

#### Get Party

Endpoint:

```text
GET /api/v1/parties/:id
```

where `:id` is the ID of the party you want to retrieve

Responses:

- `200 OK`

  Party was retrieved successfully.

  Example:

  ```json
  {
    "id": 3,
    "created_at": "2020-11-20T15:06:27.277849+01:00",
    "updated_at": "2020-11-20T15:06:27.277849+01:00",
    "name": "Amazon",
  }
  ```

- `401 Unauthorized`

  The provided token is not valid.

- `403 Forbidden`

  The party with the specified ID does not belong to the current user.

- `404 Not Found`

  The party with the specified ID does not exist.

#### Update Party

Endpoint:

```text
PATCH /api/v1/parties/:id
```

where `:id` is the ID of the party you want to update

Request payload:

```json5
{
  "name": "Rewe"
}
```

Responses:

- `200 OK`

  Party was updated successfully.

  Example:

  ```json
  {
    "id": 2,
    "created_at": "2020-11-20T15:06:27.277849+01:00",
    "updated_at": "2020-11-20T21:57:34.277849+01:00",
    "name": "Rewe"
  }
  ```

- `400 Bad Request`

  Somethinig went wrong when processing the request. Either empty request body or malformed request body.

- `401 Unauthorized`

  The provided token is not valid.

- `403 Forbidden`

  The party with the specified ID does not belong to the current user.

- `404 Not Found`

  The party with the specified ID does not exist.

- `409 Conflict`

  A party with the same name belonging to the same user already exists.

#### Delete Party

Endpoint:

```text
DELETE /api/v1/parties/:id
```

where `:id` is the ID of the party you want to delete

Responses:

- `204 No Content`

  Party was deleted successfully.

- `401 Unauthorized`

  The provided token is not valid.

- `403 Forbidden`

  The party with the specified ID does not belong to the current user.

- `404 Not Found`

  The party with the specified ID does not exist.

#### List Parties

Lists all parties which belong to the currently logged-in user.

Endpoint:

```text
GET /api/v1/parties
```

Responses:

- `200 OK`

  Parties were retrieved successfully.

  Example:

  ```json
  {
    "count": 2,
    "entries": [
      {
        "id": 2,
        "created_at": "2020-11-20T15:06:27.277849+01:00",
        "updated_at": "2020-11-20T15:06:27.277849+01:00",
        "name": "Amazon",
      },
      {
        "id": 5,
        "created_at": "2020-11-20T15:11:44.906804+01:00",
        "updated_at": "2020-11-20T15:12:19.46906+01:00",
        "name": "Rewe",
      }
    ]
  }
  ```

- `401 Unauthorized`

  The provided token is not valid.

- `404 Not Found`

  The party with the specified ID does not exist.

#### List Transactions by Party

Lists all transactions which are associated with the party with the specified ID.

Endpoint:

```text
GET /api/v1/parties/:id/transactions
```

Responses:

- `200 OK`

  Transactions were retrieved successfully.

  Example:

  ```json
  {
    "count": 1,
    "entries": [
      {
        "id": 4,
        "wallet_id": 2,
        "party_id": 2,
        "created_at": "2020-11-20T17:06:47.329695+01:00",
        "updated_at": "2020-11-20T17:06:47.329695+01:00",
        "timestamp": "2020-11-20T17:06:47.327376+01:00",
        "amount": "29.99",
        "description": "groceries"
      },
      {
        "id": 10,
        "wallet_id": 4,
        "party_id": 2,
        "created_at": "2020-11-20T17:06:47.329695+01:00",
        "updated_at": "2020-11-20T17:06:47.329695+01:00",
        "timestamp": "2020-11-20T17:06:47.327376+01:00",
        "amount": "6.99",
        "description": "cigarettes"
      }
    ]
  }
  ```

- `401 Unauthorized`

  The provided token is not valid.

- `403 Forbidden`

  The party with the specified ID does not belong to the current user.

### Transactions

A transaction is either an income (when the amount is positive) or an expense (when the amount is negative. Each transaction belongs to a specific user and is associated with a wallet and a party.)

All routes are protected and require the following header with a valid authentication token (can be obtained from [Login](#login)):

```text
Authorization: Bearer <token>
```

#### Create Transaction

Endpoint:

```text
POST /api/v1/transactions
```

Request payload:

```json5
{
  "amount": 15.50,
  "description": "Christmas decorations",
  "wallet_id": 2,
  "party_id": 2
}
```

Responses:

- `201 Created`

  Transaction was created successfully.

  Example:

  ```json
  {
    "id": 3,
    "wallet_id": 2,
    "party_id": 2,
    "created_at": "2020-11-20T15:06:27.277849+01:00",
    "updated_at": "2020-11-20T15:06:27.277849+01:00",
    "timestamp": "2020-11-20T15:06:27.277849+01:00",
    "amount": 15.50,
    "description": "Christmas decorations",
  }
  ```

- `400 Bad Request`

  Somethinig went wrong when processing the request. Either empty request body, malformed request body, missing amount or an amount of 0, missing/invalid/non-existent wallet ID, missing/invalid/non-existent party ID.

- `401 Unauthorized`

  The provided token is not valid.

- `403 Forbidden`

  The provided wallet ID or party ID does not belong to the current user.

#### Get Transaction

Endpoint:

```text
GET /api/v1/transactions/:id
```

where `:id` is the ID of the transaction you want to retrieve

Responses:

- `200 OK`

  Transaction was retrieved successfully.

  Example:

  ```json
  {
    "id": 3,
    "wallet_id": 2,
    "party_id": 2,
    "created_at": "2020-11-20T15:06:27.277849+01:00",
    "updated_at": "2020-11-20T15:06:27.277849+01:00",
    "timestamp": "2020-11-20T15:06:27.277849+01:00",
    "amount": 15.50,
    "description": "Christmas decorations",
  }
  ```

- `401 Unauthorized`

  The provided token is not valid.

- `403 Forbidden`

  The transaction with the specified ID does not belong to the current user.

- `404 Not Found`

  The transaction with the specified ID does not exist.

#### Update Transaction

Endpoint:

```text
PATCH /api/v1/transactions/:id
```

where `:id` is the ID of the transaction you want to update

Request payload:

```json5
{
  "wallet_id": 3,                                   // optional
  "party_id": 6,                                    // optional
  "timestamp": "2020-11-20T15:06:27.277849+01:00",  // optional
  "amount": 25.50,                                  // optional
  "description": "Birthday decorations",            // optional
}
```

Responses:

- `200 OK`

  Transaction was updated successfully.

  Example:

  ```json
  {
    "id": 3,
    "wallet_id": 3,
    "party_id": 6,
    "created_at": "2020-11-20T15:06:27.277849+01:00",
    "updated_at": "2020-11-20T15:06:27.277849+01:00",
    "timestamp": "2020-11-20T15:06:27.277849+01:00",
    "amount": 25.50,
    "description": "Birthday decorations",
  }
  ```

- `400 Bad Request`

  Somethinig went wrong when processing the request. Either empty request body, malformed request body, invalid/non-existent wallet ID, invalid/non-existent party ID.

- `401 Unauthorized`

  The provided token is not valid.

- `403 Forbidden`

  The provided wallet ID or party ID does not belong to the current user.

- `404 Not Found`

  The transaction with the specified ID does not exist.

#### Delete Transaction

Endpoint:

```text
DELETE /api/v1/transactions/:id
```

where `:id` is the ID of the transaction you want to delete

Responses:

- `204 No Content`

  Transaction was deleted successfully.

- `401 Unauthorized`

  The provided token is not valid.

- `403 Forbidden`

  The transaction with the specified ID does not belong to the current user.

- `404 Not Found`

  The transaction with the specified ID does not exist.

#### List all Transactions

Lists all transactions.

Endpoint:

```text
GET /api/v1/transactions
```

Responses:

- `200 OK`

  Transactions were retrieved successfully.

  Example:

  ```json
  {
    "count": 1,
    "entries": [
      {
        "id": 4,
        "wallet_id": 2,
        "party_id": 5,
        "created_at": "2020-11-20T17:06:47.329695+01:00",
        "updated_at": "2020-11-20T17:06:47.329695+01:00",
        "timestamp": "2020-11-20T17:06:47.327376+01:00",
        "amount": "29.99",
        "description": "groceries"
      },
      {
        "id": 10,
        "wallet_id": 4,
        "party_id": 2,
        "created_at": "2020-11-20T17:06:47.329695+01:00",
        "updated_at": "2020-11-20T17:06:47.329695+01:00",
        "timestamp": "2020-11-20T17:06:47.327376+01:00",
        "amount": "6.99",
        "description": "cigarettes"
      }
    ]
  }
  ```

- `401 Unauthorized`

  The provided token is not valid.

## Contributors

@desi-belokonska and @sanevillain have pair-programmed the entire project together

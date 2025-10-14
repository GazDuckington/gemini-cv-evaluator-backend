# CV evaluation backend utilizing __gemini-2.0.flash__

## Requirements

ensure the following is installed in your machine:

- [Go](https://go.dev)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Migrate](https://github.com/golang-migrate/migrate)
- [Make](https://cmake.org/)

## Getting Started

1. install go dependencies

```sh
go mod tidy
```

2. Setup environment

```sh
cp .env.example .env
```

then make sure the .env values are filled correctly.

2. Run databases

```sh
docker compose up -d --build
```
3. Migrate Data

```sh
make migrate-up
```
this will also setup 2 users account in addition to generating the required tables.

- user account:

```json
{
    "username": "user@example.com",
    "password": "123456"
}
```

- admin account:
```json
{
    "username": "admin@example.com",
    "password": "123456"
}
```

only difference is `admin` can get request to `{{host}}/users`

## Usage

to use the backend, note that default host is `http://localhost:8008`

1. login via either admin or user (see above)

```sh
POST {{host}}/login
{
    "username": "admin@example.com",
    "password": "123456"
}
```

2. submit your cv

```sh
POST {{host}}/cv
FormData:
{
"file": File(pdf),
"title": <Job title>
}
```
take the returned `id` field value

3. evaluate your cv

```sh
POST {{host}}/cv/<id>
```

it will trigger a background response that will return status, on status done the evaluation will be returned.

4. to check evaluation status

```sh
GET {{host}}/cv/status/<id>
```

5. to see evaluation response

```sh
GET {{host}}/cv/result/<id>
```

## RestAPI documentation

i use [Insomnia](https://app.insomnia.rest) as my rest client, but i have exported the collection as *HAR* file, any HTTP Client that supports *HAR* should be able to import said collection.

```sh./Insomnia_2025-10-14.har```

## Technologies

- [Gin](https://gin-gonic.com)
- [Go-Jwt](https://github.com/golang-jwt/jwt)
- [Gorm](https://gorm.io) (Postgres)
- [Logrus](https://github.com/sirupsen/logrus)
- [Godotenv](https://github.com/joho/godotenv)
- [go-genai](https://github.com/gogleapis/go-genai)
- [go-qdrant](https://github.com/qdrant/go-client)

## Databases

- postgresql 18
- qdrant latest
- minio RELEASE.2025-09-07T16-13-09Z-cpuv1

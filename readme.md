<a name="readme-top"></a>

<!-- PROJECT SHIELDS -->

<div align="center">

[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]

<!-- PROJECT LOGO -->

<br />

![Logo](https://media1.tenor.com/m/aeXj7WEQzTMAAAAd/chainsaw-chainsaw-man.gif)

  <h3 align="center">GOCRUD</h3>

  <p align="center">
    A simple Golang CRUD API
  </p>

[<img src="https://forthebadge.com/images/badges/built-with-love.svg" href="https://gorm.io" alt="gorm" width="30%">][built-with-heart-url]

</div>

<!-- ABOUT THE PROJECT -->

## About The Project

GOCRUD is a simple CRUD (Create, Read, Update, Delete) API developed in Golang. It provides basic functionalities for managing notes.

### Structure Project

```
├── assets
│   └── database-schema.png
├── cmd
│   └── main.go
├── docker-compose.yml
├── Dockerfile
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── internal
│   ├── adapter
│   │   ├── http
│   │   │   ├── handler
│   │   │   │   ├── note.go
│   │   │   │   ├── note_test.go
│   │   │   │   ├── user.go
│   │   │   │   └── user_test.go
│   │   │   ├── middleware
│   │   │   │   └── auth.go
│   │   │   └── route
│   │   │       ├── auth.go
│   │   │       ├── note.go
│   │   │       ├── user.go
│   │   │       └── welcome.go
│   │   └── repository
│   │       ├── note.go
│   │       └── user.go
│   ├── config
│   │   ├── config.go
│   │   ├── fiber.go
│   │   ├── gorm.go
│   │   └── gorm_test.go
│   ├── core
│   │   ├── domain
│   │   │   ├── claims.go
│   │   │   ├── note.go
│   │   │   ├── response.go
│   │   │   └── user.go
│   │   ├── port
│   │   │   ├── middleware.go
│   │   │   ├── note.go
│   │   │   └── user.go
│   │   └── service
│   │       ├── note.go
│   │       ├── note_test.go
│   │       ├── user.go
│   │       └── user_test.go
│   ├── mocks
│   │   ├── NoteRepository.go
│   │   ├── NoteService.go
│   │   ├── UserRepository.go
│   │   └── UserService.go
│   └── util
│       ├── bcrypt.go
│       ├── jwt.go
│       └── validator.go
└── Makefile
```

### Database Schema

![database-schema][database-schema]

### Built With

- Go
- Fiber
- GORM
- PostgreSQL

<!-- GETTING STARTED -->

## Getting Started

### Installation

Clone the repository

```bash
git clone https://github.com/shironxn/gocrud
```

### Backend

1. Navigate to the folder

   ```bash
   cd gocrud/backend
   ```

2. Copy or rename .env.example to .env

   ```bash
   cp .env.example .env
   ```

3. Check available make commands

   ```bash
   make help
   ```

#### Running Application

Run without docker

```bash
make run
```

Run with docker

```bash
make docker-up
```

Run test

```bash
make test
```

### Frontend

1. Navigate to the folder

    ```bash
    cd gocrud/frontend
    ```

2. Copy or rename .env.example to .env

    ```bash
    cp .env.example .env
    ```

3. Install package

    ```bash
    npm install
    ```    

4. Run
    ```bash
    npm run dev
    ```

<!-- API Docs -->

## API Docs

This API documentation is using [Swagger](https://swagger.io).
You can see it at http://localhost:3000/api/v1/docs

<!-- MARKDOWN LINKS & IMAGES -->

[built-with-heart-url]: https://github.com/shironxn
[contributors-shield]: https://img.shields.io/github/contributors/shironxn/gocrud.svg?style=for-the-badge
[contributors-url]: https://github.com/shironxn/blanknotes/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/shironxn/gocrud.svg?style=for-the-badge
[forks-url]: https://github.com/shironxn/blanknotes/network/members
[stars-shield]: https://img.shields.io/github/stars/shironxn/gocrud.svg?style=for-the-badge
[stars-url]: https://github.com/shironxn/blanknotes/stargazers
[issues-shield]: https://img.shields.io/github/issues/shironxn/gocrud.svg?style=for-the-badge
[issues-url]: https://github.com/shironxn/blanknotes/issues
[database-schema]: ./assets/database-schema.png
[golang-shield]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[golang-url]: https://go.dev
[fiber-url]: https://gofiber.io
[gorm-url]: https://gorm.io
[mysql-shield]: https://img.shields.io/badge/mysql-4479A1.svg?style=for-the-badge&logo=mysql&logoColor=white
[mysql-url]: https://www.mysql.com

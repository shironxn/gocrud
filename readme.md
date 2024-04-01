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

[<img src="https://forthebadge.com/images/badges/built-with-love.svg" href="https://gorm.io" alt="gorm" width="50%">][built-with-heart-url]

</div>

<!-- ABOUT THE PROJECT -->

## About The Project

GOCRUD is a simple CRUD (Create, Read, Update, Delete) API developed in Golang. It provides basic functionalities for managing notes.

![database-schema][database-schema]

### Built With

- Go
- Fiber
- GORM
- MySQL

<!-- GETTING STARTED -->

## Getting Started

### Installation

1. Clone the repository

   ```bash
   git clone https://github.com/shironxn/gocrud.git
   ```

2. Navigate to the project directory

   ```bash
   cd gocrud
   ```

3. Copy or rename .env.example to .env

   ```bash
   cp .example.env .env
   ```

4. Check available make commands

   ```bash
   make help
   ```

### Running Application

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

<!-- API Docs -->

## API Docs

This API documentation is using [Swagger](https://swagger.io).
You can see it at http://localhost:3000/api/v1/docs

<!-- MARKDOWN LINKS & IMAGES -->

[built-with-heart-url]: https://github.com/shironxn
[contributors-shield]: https://img.shields.io/github/contributors/shironxn/gocrud.svg?style=for-the-badge
[contributors-url]: https://github.com/shironxn/gocrud/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/shironxn/gocrud.svg?style=for-the-badge
[forks-url]: https://github.com/shironxn/gocrud/network/members
[stars-shield]: https://img.shields.io/github/stars/shironxn/gocrud.svg?style=for-the-badge
[stars-url]: https://github.com/shironxn/gocrud/stargazers
[issues-shield]: https://img.shields.io/github/issues/shironxn/gocrud.svg?style=for-the-badge
[issues-url]: https://github.com/shironxn/gocrud/issues
[database-schema]: ./assets/database-schema.png
[golang-shield]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[golang-url]: https://go.dev
[fiber-url]: https://gofiber.io
[gorm-url]: https://gorm.io
[mysql-shield]: https://img.shields.io/badge/mysql-4479A1.svg?style=for-the-badge&logo=mysql&logoColor=white
[mysql-url]: https://www.mysql.com

## About Visor

Visor is a quick and simple way to get your next Laravel project running locally with Docker on Mac OS.

Install Visor globally with composer:

```bash
composer global require thebrubaker/laravel-visor:dev-linux
```

To use visor to spin up a new Laravel application, simply run `visor up`

Tip: If you are already familiar with Docker or Docker Compose, this tool is likely only useful for convenience.

## Visor Commands

```bash
visor up # spin up application
visor down # spin down application
visor migrate # run migrations
visor tinker # jump into your php container with a bash shell
visor compose ...args # run docker-compose commands
```

## Requirements

This tool requires Docker and should run on Mac OS. If you want to install the linux binary, you can require `thebrubaker/laravel-visor:dev-linux`

## Visor Up

```
👉 running composer install...
👉 spinning up services...
👉 running migrations...

💪 run `visor down` to spin down your application and services
💪 run `visor tinker` to jump into a php container

👌 Applicaton available at http://localhost:8080
👌 Database available at mysql://root:secret@127.0.0.1:3306/laravel_visor
```

## Visor Down

```
👉 spinning down services...
```

## Visor Init

```
It looks like this is your first time running Visor. Proceed with visor init? (Y/n)

👌 Docker is installed
👌 Created .visor directory and added to .gitignore

👉 downloading containers for php 7.4, redis 4.0 and mysql 5.7...

👌 Visor init success!
```

## Visor Errors

```
💥 Docker is not installed on this machine

Visor requires Docker to continue. Have you installed Docker on your machine?

💥 Unable to acquire DB_PORT=3306 from your .env file.

Visor exposes a mysql database on the port listed in your .env file. Please update DB_PORT to one available on your machine and try again.

💥 A docker-compose.yaml file already exists for this project.

Visor is a simple wrapper for docker-compose. Do you want us to backup your docker-compose config and replace it with Visor's config?

💥 Unable to access your application at http://localhost:8080

Visor attempted to spin up your application on port 8080. Was this port already taken?

💥 Unable to access your database at mysql://root:secret@127.0.0.1:3306/laravel_visor

Visor attempted to spin up your database on port 3306. Was this port already taken?
```

## Visor Build

Tasks

- Build docker image
- Upload docker image

## Visor Deploy

Tasks

- Build docker image w/ tag
- Upload docker image w/ tag
- Trigger update to Cloud Run

## Goals

- Run a local Laravel application from scratch
- Set up google cloud run (free tier)
- Deploy to google cloud run
- Access via ssh into a production server

```
visor cloud init
visor cloud build
visor cloud deploy
```

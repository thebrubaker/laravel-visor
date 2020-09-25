## About Visor

Visor is a dead simple way to manage your next Laravel project locally and into production. This tool optimizes for a dead-simple setup and sane-defaults.

-   One command to run a local Laravel application from scratch
-   One command to set up google cloud run (free tier)
-   One command to deploy to production
-   One command to ssh into a production server

```bash
visor up
visor down
visor cloud init
visor cloud build
visor cloud deploy
visor cloud exec
visor cloud info
```

### Requirements

This tool requires Docker and runs on Mac / Linux / Windows w/ WSL2.

## First Time Setup

It looks like this is your first time running Visor. Proceed with visor init? (Y/n)

👌 Docker is installed
👌 Created .visor directory and added to .gitignore

👉 downloading containers for php 7.4, redis 4.0 and mysql 5.7...

👌 Visor init success!

## Errors

💥 Docker is not installed on this machine

## Visor Up

👉 spinning up services...
👉 running composer install...

💪 run `visor migrate` to migrate your database
💪 run `visor down` to spin down your application and services
💪 run `visor tinker` to jump into the php container

👌 Applicaton available at http://localhost:8080
👌 Database available at mysql://root:secret@127.0.0.1:3306/laravel_visor

## Visor Down

Tasks

-   Check if docker-compose file exists
-   Run docker-compose down

## Visor Init

Tasks

-   Check if gcloud SDK is installed
-   Ask user for services to be created
-   Create Project
-   Create Cloud Run
-   Create Compute Instance
-   Create Storage Bucket
-   Create .env.prod file in repo
-   Add .env.prod to .gitignore

## Visor Build

Tasks

-   Build docker image
-   Upload docker image

## Visor Deploy

Tasks

-   Build docker image w/ tag
-   Upload docker image w/ tag
-   Trigger update to Cloud Run

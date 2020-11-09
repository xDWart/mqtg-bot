# mqtg-bot

[![Build Status](https://travis-ci.com/xDWart/mqtg-bot.svg?branch=master)](https://travis-ci.com/xDWart/mqtg-bot)
[![Go Report Card](https://goreportcard.com/badge/github.com/xDWart/mqtg-bot)](https://goreportcard.com/report/github.com/xDWart/mqtg-bot)
[![Version](https://img.shields.io/github/go-mod/go-version/xDWart/mqtg-bot)](go.mod)
[![License](https://img.shields.io/github/license/xDWart/mqtg-bot)](LICENSE)

![Gopher Bot](https://github.com/xDWart/mqtg-bot/raw/master/assets/kdpv.jpg)

## Articles

- [Flexible MQTT-telegram-bot for IoT](https://xd-wart.medium.com/flexible-mqtt-telegram-bot-for-iot-70d567edfb2e) (En, Medium)
- [(Не)очередной MQTT-телеграм-бот для IoT](https://habr.com/ru/post/526672/) (Ru, Habr)

## Introduction

mqtg-bot is an easy-to-configure for your needs MQTT client Telegram bot. Without programming knowledge you can configure the bot to send various commands (turn on the light, open the garage door, etc.) or request any information (the temperature in the house, the state of the heating system, etc.) or receive frames from security camera. In general the functionality of this solution is very rich.

![Edit buttons menu](https://github.com/xDWart/mqtg-bot/raw/master/assets/edit_buttons_menu.jpg)

![Temp and Humidity](https://github.com/xDWart/mqtg-bot/raw/master/assets/temp_and_humidity.jpg)

![Take a picture](https://github.com/xDWart/mqtg-bot/raw/master/assets/take_a_picture.jpg)

## Features

- [x] Connecting to MQTT broker
    - [x] tcp / ssl / ws / wss
- [x] Subscribing to a topic:
    - [x] Selectable QoS/Retained
    - [x] Text/Image data types
    - [x] Pre/post value displaying text
    - [x] Storing data into DB
    - [ ] Data storage management
    - [ ] Publish action on receiving
    - [ ] Parse data from message by regexp
    - [ ] Voice data type
- [x] Publishing to a topic:
    - [x] Selectable QoS/Retained
    - [x] Text/Image data types
    - [ ] Voice data type
- [x] Customized users buttons menu:
    - [x] Folders
    - [x] Single-value buttons
    - [x] Toggle buttons
    - [x] Multi-value buttons
    - [x] Print last subscription value
    - [x] Draw charts
- [ ] Your great idea ([create a proposal in issues](https://github.com/xDWart/mqtg-bot/issues/new/choose))

## Usage

You can run the bot on your Raspberry Pi home server or free Heroku dyno.

Clone this repository:
```sh
git clone https://github.com/xDWart/mqtg-bot
```

Message [@BotFather](https://telegram.me/BotFather) `/newbot` command to create a bot and get its HTTP API access token.

#### Required environment variables

- `TELEGRAM_BOT_TOKEN` (HTTP API access token)
- `DATABASE_URL` **or** `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`

Notes: 
1. `DATABASE_URL` must have the following format: `postgres://user:password@host:port/db`
1. You can create the `.env` file in the root of the project and insert your key/value environment variable pairs in the following format of `KEY=VALUE`

#### Local running

If you don't have Postgres, the easiest way to get it is to run it under Docker:

```sh
docker run --name some-postgres -e POSTGRES_PASSWORD=... -p 5432:5432 -d postgres
```

You can run mqtg-bot with environment variables:

```sh
TELEGRAM_BOT_TOKEN=... POSTGRES_PASSWORD=... go run main.go
```

or if you've already created the `.env` file:

```sh
go run main.go
```

#### Running in Docker container

```sh
docker run -e TELEGRAM_BOT_TOKEN=... -e DATABASE_URL=... --network=host owart/mqtg-bot
```

#### Running both mqtg-bot and Postgres with Docker Compose

Create a `.env` file in the root of the project and insert your key/value environment variable pairs in the following format of `KEY=VALUE`.

```sh
TELEGRAM_BOT_TOKEN=*BOT_ACCESS_TOKEN*
POSTGRES_PASSWORD=password
```

Run docker-compose:

```sh
docker-compose up --build
```

#### Heroku running

You will need Heroku CLI

```sh
# login into Heroku
heroku login

# create a new app
heroku create

# add your new app into git remotes
heroku git:remote -a *YOUR_APP_NAME*

# add TELEGRAM_BOT_TOKEN environment
heroku config:set TELEGRAM_BOT_TOKEN=*BOT_ACCESS_TOKEN*

# attach Postgres add-on
heroku addons:create heroku-postgresql:hobby-dev 

# push master branch to Heroku
git push heroku master

# scale up your app
heroku ps:scale worker=1
```

Then just message `/start` to your bot and follow the instructions to configure it.

## Try bot

[@mqtg_bot](https://telegram.me/mqtg_bot)  
Message `/start` to him and configure connection to your MQTT broker.

## Contribution

- Do you have an idea to improve mqtg-bot? -> [Create an issue](https://github.com/xDWart/mqtg-bot/issues/new/choose).     
- Have you discovered a bug? -> [Create an issue](https://github.com/xDWart/mqtg-bot/issues/new/choose).   
- Have you already coded something for mqtg-bot? -> [Create a pull request](https://github.com/xDWart/mqtg-bot/compare).   

## Licence

- mqtg-bot is licensed under the MIT License.   
- See [LICENSE](LICENSE) for the full license text.   
- Copyright (c) Anatoliy Bezgubenko
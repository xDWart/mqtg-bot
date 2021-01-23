# mqtg-bot

[![Build Status](https://github.com/xDWart/mqtg-bot/workflows/Build/badge.svg)](https://github.com/xDWart/mqtg-bot/actions?query=workflow%3ABuild)
[![Go Report Card](https://goreportcard.com/badge/github.com/xDWart/mqtg-bot)](https://goreportcard.com/report/github.com/xDWart/mqtg-bot)
[![Version](https://img.shields.io/github/go-mod/go-version/xDWart/mqtg-bot)](go.mod)
[![License](https://img.shields.io/github/license/xDWart/mqtg-bot)](LICENSE)

![Gopher Bot](https://github.com/xDWart/mqtg-bot/raw/master/assets/kdpv.jpg)

## Articles

- [Development of the Open-Source Telegram Bot for MQTT IoT](https://dzone.com/articles/development-of-the-open-source-telegram-bot-for-mq) (En, DZone.com)
- [(Не)очередной MQTT-телеграм-бот для IoT](https://habr.com/ru/post/526672/) (Ru, Habr)
- [Разработка NoCode решения для интернета вещей](https://youtu.be/Ja2AjOAGnnY) (Ru, YouTube)

## Introduction

mqtg-bot is an easy-to-configure for your needs MQTT client Telegram bot. Without programming knowledge you can configure the bot to send various commands (turn on the light, open the garage door, etc.) or request any information (the temperature in the house, the state of the heating system, etc.) or receive frames from security camera. In general the functionality of this solution is very rich.

![Edit buttons menu](https://github.com/xDWart/mqtg-bot/raw/master/assets/edit_buttons_menu.jpg)

![Temp and Humidity](https://github.com/xDWart/mqtg-bot/raw/master/assets/temp_and_humidity.jpg)

![Take a picture](https://github.com/xDWart/mqtg-bot/raw/master/assets/take_a_picture.jpg)

## Features

- [x] Connecting to MQTT broker
    - [x] tcp / ssl / ws / wss
- [x] Supported databases
    - [x] Postgres
    - [x] SQLite
- [x] Subscribing to a topic:
    - [x] Selectable QoS/Retained
    - [x] Text/Image data types
    - [x] Pre/post value displaying text
    - [x] Storing data into DB
    - [x] Parse data by JsonPath expressions
    - [ ] Data storage management
    - [ ] Publish action on receiving
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

Message [@BotFather](https://telegram.me/BotFather) `/newbot` command to create a bot and get his HTTP API access token.

#### Environment variables

- `TELEGRAM_BOT_TOKEN` - bot HTTP API access token, required
- `DATABASE_URL` - Postgres connection string in the following format: `postgres://user:password@host:port/db`
- `SQLITE_PATH` - path to SQLite database

Notes: 
1. Only `TELEGRAM_BOT_TOKEN` env is required
1. If `DATABASE_URL` env is omitted, or a Postgres connection error occurred, SQLite will be used
1. If `SQLITE_PATH` env is omitted, `mqtg.db` will be used by default as a SQLite database
1. You can create the `.env` file in the root of the project and insert your key/value environment variable pairs in the following format of `KEY=VALUE`

#### Local running

You can run mqtg-bot with environment variables:

```sh
TELEGRAM_BOT_TOKEN=... go run main.go
```

or if you've already created the `.env` file:

```sh
go run main.go
```

#### Running in Docker container

```sh
docker run -e TELEGRAM_BOT_TOKEN=... -e DATABASE_URL=... --network=host owart/mqtg-bot
```

#### Heroku running

You will need Heroku CLI

```sh
# login into Heroku
heroku login

# create a new app
heroku create *YOUR_APP_NAME*

# add your new app into git remotes
heroku git:remote -a *YOUR_APP_NAME*

# add TELEGRAM_BOT_TOKEN environment
heroku config:set TELEGRAM_BOT_TOKEN=*BOT_ACCESS_TOKEN*

# set version of Go
heroku config:set GOVERSION=go1.15

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

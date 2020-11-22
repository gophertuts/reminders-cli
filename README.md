# Reminders CLI app

## Overview

In this project we'll build a **multi component CLI application**
which consists of mainly 3 parts: a `CLI client`, a `backend API` server
and a `notifier service`.

The CLI client will take the input from command line
and pass it to the backend API through its HTTP client.

The backend API server is an **HTTP server** which
has all the needed endpoints for **CRUD operations** on reminders.
It also has 2 background running workers: **background saver**
& **background notifier**. Correspondingly saving the in-memory
data to the disk and notifying un-completed reminders.

The backend API server also communicates with the notifier service
through its own HTTP client.

Speaking of the database layer, we'll be creating our own **file database storage**
with some optimized mechanism for this type of application.

Ta-dah 🥳 🚀

<img src="https://github.com/gophertuts/reminders-cli/raw/master/cli-demo.gif?sanitize=true"/>

## Medium article 📖

- [Reminders CLI in Go](https://www.youtube.com/c/GopherTuts)

## YouTube tutorials 🎥

- [Reminders CLI in Go #1 - Project setup & bare bones](https://youtu.be/-9CbX2MncZg) - [[Download Code]](https://github.com/gophertuts/reminders-cli/raw/master/zips/reminders-cli-1.tar.gz)
- [Reminders CLI in Go #2 - Notifier Service](https://youtu.be/rlsnqlSjUOc) - [[Download Code]](https://github.com/gophertuts/reminders-cli/raw/master/zips/reminders-cli-2.tar.gz)
- [Reminders CLI in Go #3 - CLI Basics](https://youtu.be/PbKCvQuAPIQ) - [[Download Code]](https://github.com/gophertuts/reminders-cli/raw/master/zips/reminders-cli-3.tar.gz)
- [Reminders CLI in Go #4 - Command Switch - Part 1](https://youtu.be/Vz2dBY_hAkw) - [[Download Code]](https://github.com/gophertuts/reminders-cli/raw/master/zips/reminders-cli-4.tar.gz)
- [Reminders CLI in Go #5 - Command Switch - Part 2](https://youtu.be/Py10z9-61JQ) - [[Download Code]](https://github.com/gophertuts/reminders-cli/raw/master/zips/reminders-cli-5.tar.gz)


## Requirements 🤓

- [Go](https://golang.org/doc/install)
- [Node.js](https://nodejs.org/en/download/)

In this tutorial we'll be writing a little bit of [Node.js](https://nodejs.org/en/download/)
aka the `Notifier Service` because it's the fastest
cross platform OS notification system available for us.

We'll also be using [Yarn](https://yarnpkg.com/lang/en/docs/install/) package manager for this application.

And that's all on the JavaScript (Node.js) side.
The rest is pure `Go code` also **without** any **third party packages**
meaning we'll write absolutely everything from scratch.

## Components 🧩

- **CLI Client**
- **HTTP client** for communicating with the Backend API
- **Backend API**
- **HTTP client** for communicating with the Notifier service
- **Notifier** service
- Background **Saver worker**
- Background **Notifier worker**
- JSON file **Database** (`db.json`)
- **Database config** file (`.db.config.json`)

## CLI Client

#### Features

- `create` a reminder
- `edit` a reminder
- `fetch` a list of reminders
- `delete` a list of reminders

***Note:*** Only works if Backend API is up & running

## Backend API

#### Features

- Does CRUD operations with incoming data from CLI client
- Runs Background Saver worker, which saves in-memory data
- Runs Background Notifier worker, which notifies un-completed reminders
- It can work without the Notifier service, and will keep
retrying unsent notifications until Notifier service is up
- On backend API shutdown all the in-memory data is saved

#### Endpoints

- `GET /health`                 - responds with 200 when server is up & running 
- `POST /reminders/create`      - creates a new reminder and saves it to DB
- `PUT /reminders/edit`         - updates a reminder and saves it to DB (if duration is updated, notification is resent)
- `POST /reminders/fetch`       - fetches a list of reminders from DB
- `DELETE /reminders/delete`    - deletes a list of reminders from DB

## Background Saver

#### Features

- Saves in-memory reminders to the disk (`db.json`)
- Saves db config to the disk (`.db.config.json`)

## Background Notifier

#### Features

- Pushes un-completed reminders to the Notifier service

## Notifier Service

#### Features

- Sends OS notifications

#### Endpoints

- `GET /health`                 - responds with 200 when server is up & running
- `POST /notify`                - sends OS notification and retry response

## File DB

#### Features

- Records are saved inside `db.json` file
- Has a db config file (`.db.config.json`)
- Has an auto increment ID generator

## Installation ⚙

Before running any command or trying to compile the programs
make sure you first have all the needed dependencies installed:

- [Golang](https://golang.org/doc/install)
- [GoLint](https://github.com/golang/lint)
- [Node.js](https://nodejs.org/en/download/)
- [Node.js Ubuntu](https://tecadmin.net/install-latest-nodejs-npm-on-ubuntu/)
- [Yarn](https://yarnpkg.com/lang/en/docs/install/)
- [GitBash - WINDOWS ONLY](https://git-scm.com/download/win)
- [Cygwin - WINDOWS ONLY](https://www.cygwin.com/)
- [Make](https://sourceforge.net/projects/ezwinports/files/make-4.2.1-without-guile-w32-bin.zip/download)

###### Configure `make` (WINDOWS ONLY):

***Note:*** Make sure you have [GitBash](https://git-scm.com/download/win) installed
before proceeding.

1. Download the [Make](https://sourceforge.net/projects/ezwinports/files/make-4.2.1-without-guile-w32-bin.zip/download)
executable

2. Extract the contents form the zip

3. Place the `bin/make.exe` inside `C:\Program Files\Git\mingw64\bin`

4. If you're using **Goland** update your SHELL
`Ctrl` + `Alt` + `S` `-->` `Tools` `-->` `Terminal` `-->` `Shell Path` `-->` `"C:\Program Files\Git\bin\sh.exe" --login -i`

5. Restart Goland IDE

For more info refer to [GitBash - CygWin](https://gist.github.com/evanwill/0207876c3243bbb6863e65ec5dc3f058)

---

## Run 🎮

#### `make` commands

```bash
# builds client & server binaries, formats & lints the code
make

# builds the client binary
make client

# builds the server binary
make server

# formats the entire code base
make fmt

# lints the entire code base
make lint

# checks the entire code base for code issues
make vet
```

#### `server` flags

```bash
# display a helpful message of all available flags for the server binary
./bin/server --help

# runs the backend http server on the specified address
# --backend flag needs to be provided to ./bin/client if address != :8080
./bin/server --addr=":9090"

# runs the http backend server with a different path to the database
./bin/server --db="/path/to/db.json"

# runs the http backend server with a different path to the database config
./bin/server --db-cfg="/path/to/.db.config.json"

# runs the http backend server with a different notifier service url
./bin/server --notifier="http://localhost:8989"
```

#### `client` commands & flags

```bash
# displays a helpful message about all the commands and flags available
./bin/client --help

# runs CLI client with a different backend api url
./bin/client --backend="http://localhost:7777"

# creates a new reminder which will be notified after 3 minutes
./bin/client create --title="Some title" --message="Some msg!" --duration=3m

# edits the reminder with id: 13
# note: if the duration is edited, the reminder gets notified again
./bin/client edit --id=13 --title="Another title" --message="Another msg!"

# fetches a list of reminders with the following ids
./bin/client fetch --id=1 --id=3 --id=6

# deleted the reminders with the following ids
./bin/client delete --id=2 --id=4
```

---

***Note:*** Before using `./bin/client` binary,
make sure to have `/bin/server` and `notifier/notifier.js` up & running

**1st terminal**
```bash
node notfier/notifier.js
```

**2nd terminal**
```bash
./bin/server
```

**3rd terminal**
```bash
./bin/client ...
```

## Resources 💎

- [Handler](https://golang.org/pkg/net/http/#Handler)

## Feedback 🐧

[SteveHook TypeForm](https://feedback.gophertuts.com)

## Community 💬

[SteveHook Discord](https://discord.gg/tprewQu)

---

**Enjoy** 🚀🚀🚀

<img src="https://github.com/gophertuts/go-basics/raw/master/gophertuts.svg?sanitize=true" width="50px"/>

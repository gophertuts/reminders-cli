# Reminders CLI app

## Overview

<img src="https://github.com/gophertuts/reminders-cli/raw/master/cli-demo.gif?sanitize=true"/>

## Requirements

- Node.js
- Go

## Components

- CLI interface
- HTTP client for communicating with Backend API
- Backend API
- HTTP client for communicating with Notifier service
- Notifier service
- Background Saver worker
- Background Notifier worker

## CLI Client

#### Features

- `CREATE` reminder
- `EDIT` reminder
- `FETCH` a list of reminders
- `DELETE` a list of reminders
- Only works if Backend API is up and running

## Backend API

#### Features

- Does CRUD operations with incoming data from CLI client
- Runs Background Saver worker, which saves in-memory data
- Runs Background Notifier worker, which notifies un-completed reminders
- It can work without the Notifier service, and will keep
retrying unsent notifications until Notifier service is up

#### Endpoints

- `GET /health`                 - responds with 200 when server is up & running 
- `POST /reminders/create`      - creates a new reminder and saves it to DB
- `PUT /reminders/edit`         - updates a reminder and saves it to DB (if duration is updated, notification is resent)
- `POST /reminders/fetch`       - fetches a list of reminders from DB
- `DELETE /reminders/delete`    - deletes a list of reminders from DB

## Notifier

#### Features

- Sends OS notifications

#### Endpoints

- `GET /health`                 - responds with 200 when server is up & running
- `POST /notify`                - sends OS notification and retry response

## Installation

Before running any command or trying to compile the programs
make sure you first have all the needed dependencies:

- [Golang](https://golang.org/doc/install)
- [Golint](https://github.com/golang/lint)
- [Node.js](https://nodejs.org/en/download/)
- [Node.js Ubuntu](https://tecadmin.net/install-latest-nodejs-npm-on-ubuntu/)
- [Yarn](https://yarnpkg.com/lang/en/docs/install/)
- [GitBash - WINDOWS ONLY](https://git-scm.com/download/win)
- [Cygwin - WINDOWS ONLY](https://www.cygwin.com/)
- [Make](https://sourceforge.net/projects/ezwinports/files/make-4.2.1-without-guile-w32-bin.zip/download)

Configure CygWin (WINDOWS ONLY):

1. Download the [Make](https://sourceforge.net/projects/ezwinports/files/make-4.2.1-without-guile-w32-bin.zip/download)
executable

2. Extract the contents form the zip

3. Place the `bin/make.exe` inside `C:\Program Files\Git\mingw64\bin`

4. If you're using Goland update your SHELL
`Ctrl` + `Alt` + `S` --> `Tools` --> `Terminal` --> `Shell Path` --> `"C:\Program Files\Git\bin\sh.exe" --login -i`

5. Restart Goland IDE

For more info refer to [GitBash - CygWin](https://gist.github.com/evanwill/0207876c3243bbb6863e65ec5dc3f058)

---

## Run

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

Before using `./bin/client` binary, make sure to have `/bin/server` and `notifier/notifier.js`
up and running

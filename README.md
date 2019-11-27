# migrate

MongoDB migration tool with minimal api.

## Features

- Azure CosmosDB ready
- Manage schema in pure JSON
- No magic, all configurations comes from standard MongoDB commands

## Commands

Export MongoDB connection as URI before running any commands.

    export URI="mongodb://username:password@example.com:27017/database?ssl=true&retrywrites=false"

### URI format

    "mongodb://<username>:<password>@<hostname>:<port>/<db>?<options>"

### Init

    migrate init

### Run migration

    migrate up

### Check status

    migrate status

### Rollback migration

We don't have such command, take care by yourself.

## How it works

The JSON schema must have the following format.

    {
      "adminCommand": {}
      "command": {}
    }

- Anything inside `adminCommand` goes to `db.admin().runCommand({})`
- `command` goes to `db.runCommand({})`

See [examples](examples) for more information.

## Development

### Requirements

Install these dependencies into your machine.

- Go 1.13+
- Docker
- docker-compose

### Install dependencies

    make install

### Run from source

    make run

### Test

    make test

### Build

Using [gox](https://github.com/mitchellh/gox) to build cross platform binaries.

    make build

# migrate

MongoDB migration tool with minimal api.

## Features

- Azure CosmosDB ready
- Manage schema in pure JSON
- No magic, all configurations comes from standard MongoDB commands

## Install

This module is currently not exposed in any open source repository, so you will need to install by yourself.

Therefore release artifact is checked into this repository.

1. Copy binary for your OS from [releases](https://sios.tech/covas/migrate/releases)
2. Move the binary to `/usr/local/bin/migrate`

if you create collections, you need to install [azure-cli](https://pypi.org/project/azure-cli/).

    pip install azure-cli

## Commands

Export MongoDB connection as URI before running any commands.

    export URI="mongodb://username:password@example.com:27017/database?ssl=true&retrywrites=false"

### URI format

    "mongodb://<username>:<password>@<hostname>:<port>/<db>?<options>"

### Init

    migrate init

### Run migration

    migrate up -d "migrations" -r "develop"

### Fix migration

マイグレーションに失敗したファイルを再マイグレーションするコマンド。

    migrate fix -f <ファイル名> -a false -r <リソースグループ>

### Revert migration pointer

DB に記録されているマイグレーションポインタをリセットするコマンド。
（DB に現在のマイグレーション状況をファイル名で保存していく仕組みなので、失敗したときはそれを巻き戻すコマンドが必要になる。）

    migrate revert -n <ファイル名>

### Find index

コレクションに記録されたインデックスを確認するコマンド（標準出力で表示される）。

    migrate index -n <コレクション名>

### Delete index

インデックスの名前を指定して削除するコマンド。

    migrate delete -n <インデックス名> -c <コレクション名>

### Rollback migration

There is no such command, take care by yourself.

## How it works

The JSON schema must have the following format.

    {
      "adminCommand": {}
      "command": {}
    }

- Anything inside `adminCommand` goes to `az cosmosdb collection **`
- `command` goes to `db.runCommand({})`

See [examples](examples-v2) for more information.

## Development

### Requirements

Install these dependencies into your machine.

- Go 1.20+
- Docker
- docker-compose

### Install dependencies

    make install

### Run from source

    make run

### Test

    make test

### Build

Build cross platform binaries.

    make build

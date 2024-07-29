# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.6.1] - 2024-07-29

### Changed

- shardKey を必須パラメータに変更。shardKey が存在しない場合にエラーを返す。

## [0.6.0] - 2024-07-26

### Changed

- adminCommand を azure cli で実行するように変更

## 0.5.3

- Changed: Azure CosmosDB for MongoDB API 仮想コアの接続先 URI に対応

## 0.5.2

- Changed: update mongo-go-driver

## 0.5.1

- Feature: When Failed run command, return exit 1

## 0.5.0

- Changed: update go version
- Changed: update go packeges

## 0.2.0

- Feature: Zipped each builds
- Bugfix: Referenced to nil object

## 0.1.0

- Cleanup: Put binaries on releases
- Feature: Ensured zero string value when adminCommand or command does not exist
- Feature: Added README
- Refactor: Flattened packages to become top-level
- Feature: Installed gox for multi platform builds
- Feature: Added docker-compose
- Feature: Init

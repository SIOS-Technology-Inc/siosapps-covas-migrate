package main

// 任意のazコマンドを実行する

import (
	"fmt"
	"os/exec"
)

// AzureCommand はAzure Cosmos DB for MongoDBのコレクションを管理するためのコマンドを作成するための構造体
type AzureCommand struct {
	Description string `json:"description"`
	Collection  string `json:"collection"`
	ShardKey    string `json:"shardKey"`
	SharedRU    bool   `json:"sharedRU"`
	AutoScale   *bool  `json:"autoScale"`
	Throughput  *int   `json:"throughput"`
}
type Action string

const (
	Create Action = "create"
	Delete Action = "delete"
	List   Action = "list"
	Show   Action = "show"
	Exists Action = "exists"
)

/*
Azure Cosmos DB for MongoDBのコレクションを管理するためのコマンドを作成する

	※コレクション作成時の注意点
	データベース共有RUからコンテナー固有RUには変更できません。
	コンテナー固有RU/s の手動スループットを設定する場合は、--throughput パラメーターを使用します。
	az cosmosdb mongodb collection create -g MyResourceGroup -a MyAccount -d MyDatabase -n MyCollection --shard "ShardingKey" --idx @indexes-file.json --throughput "400"
	コンテナー固有RU/s の自動スループットを設定する場合は、--max-throughput パラメーターを使用します。
	az cosmosdb mongodb collection create -g MyResourceGroup -a MyAccount -d MyDatabase -n MyCollection --shard "ShardingKey" --idx @indexes-file.json --max-throughput "4000"
	データベース共有RU/s のスループットを設定する場合は、--throughput　--max-throughput パラメーターを指定しません。
	az cosmosdb mongodb collection create -g MyResourceGroup -a MyAccount -d MyDatabase -n MyCollection --shard "ShardingKey" --idx @indexes-file.json
*/
func (ac *AzureCommand) CreateCommand(action Action, rg string, accountName string, dbName string) ([]string, error) {
	base := []string{"cosmosdb", "mongodb", "collection"}

	if ac.Collection == "" {
		return nil, fmt.Errorf("collection name is required")
	}
	if rg == "" {
		return nil, fmt.Errorf("resource group name is required")
	}
	if accountName == "" {
		return nil, fmt.Errorf("account name is required")
	}
	// コマンドの作成
	var args []string
	switch action {
	case Create:
		args = append(base, "create", "-g", rg, "-a", accountName, "-d", dbName, "-n", ac.Collection)
		if ac.ShardKey != "" {
			args = append(args, "--shard", ac.ShardKey)
		} else {
			return nil, fmt.Errorf("shard key is required")
		}
		// データベース共有RU/sがオフかつthroughputが指定されていない場合はエラー
		if !ac.SharedRU && ac.Throughput == nil {
			return nil, fmt.Errorf("throughput and maxThroughput cannot be specified at the same time")
		}
		if ac.SharedRU {
			break
		} else {
			if ac.AutoScale != nil && *ac.AutoScale {
				if *ac.Throughput < 4000 {
					return nil, fmt.Errorf("maxThroughput must be greater than or equal to 4000")
				}
				args = append(args, "--max-throughput", fmt.Sprint(*ac.Throughput))
			} else {
				if *ac.Throughput < 400 {
					return nil, fmt.Errorf("throughput must be greater than or equal to 400")
				}
				args = append(args, "--throughput", fmt.Sprint(*ac.Throughput))
			}
		}
	case Delete:
		args = append(base, "delete", "-g", rg, "-a", accountName, "-d", dbName, "-n", ac.Collection)
	case List:
		args = append(base, "list", "-g", rg, "-a", accountName, "-d", dbName)
	case Show:
		args = append(base, "show", "-g", rg, "-a", accountName, "-d", dbName, "-n", ac.Collection)
	case Exists:
		args = append(base, "exists", "-g", rg, "-a", accountName, "-d", dbName, "-n", ac.Collection)
	default:
		return nil, fmt.Errorf("invalid action")

	}
	return args, nil
}

// AzExcute は任意のazコマンドを実行する
func AzExcute(args []string) ([]byte, error) {
	// azコマンドを実行
	cmd := "az"
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		fmt.Println("az command failed")
	}
	return out, err
}

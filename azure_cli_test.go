package main

import (
	"reflect"
	"testing"
)

func TestCollection(t *testing.T) {
	ac := AzureCommand{
		Collection: "MyCollection",
		ShardKey:   "_id",
		SharedRU:   false,
		AutoScale:  nil,
		Throughput: nil,
	}

	rg := "MyResourceGroup"
	accountName := "MyAccount"
	dbName := "MyDatabase"

	t.Run("success - 固有RU autoscale無効", func(t *testing.T) {
		action := Create
		expectedArgs := []string{"cosmosdb", "mongodb", "collection", "create", "-g", rg, "-a", accountName, "-d", dbName, "-n", ac.Collection, "--shard", "_id", "--throughput", "400"}
		ac.ShardKey = "_id"
		ac.AutoScale = new(bool)
		*ac.AutoScale = false
		throughput := 400
		ac.Throughput = &throughput

		args, err := ac.CreateCommand(action, rg, accountName, dbName)

		if !reflect.DeepEqual(args, expectedArgs) {
			t.Errorf("expected args %v, got %v, err %v", expectedArgs, args, err)
		}
	})

	t.Run("success - 固有RU autoscale有効", func(t *testing.T) {
		action := Create
		expectedArgs := []string{"cosmosdb", "mongodb", "collection", "create", "-g", rg, "-a", accountName, "-d", dbName, "-n", ac.Collection, "--shard", "_id", "--max-throughput", "4000"}
		ac.ShardKey = "_id"
		ac.AutoScale = new(bool)
		*ac.AutoScale = true
		throughput := 4000
		ac.Throughput = &throughput

		args, err := ac.CreateCommand(action, rg, accountName, dbName)

		if !reflect.DeepEqual(args, expectedArgs) {
			t.Errorf("expected args %v, got %v, err %v", expectedArgs, args, err)
		}
	})

	t.Run("success - 共有RU", func(t *testing.T) {
		action := Create
		expectedArgs := []string{"cosmosdb", "mongodb", "collection", "create", "-g", rg, "-a", accountName, "-d", dbName, "-n", ac.Collection, "--shard", "_id"}
		ac.ShardKey = "_id"
		ac.SharedRU = true

		args, err := ac.CreateCommand(action, rg, accountName, dbName)

		if !reflect.DeepEqual(args, expectedArgs) {
			t.Errorf("expected args %v, got %v, err %v", expectedArgs, args, err)
		}
	})

	t.Run("error - collection name is required", func(t *testing.T) {
		action := Create
		ac.Collection = ""

		_, err := ac.CreateCommand(action, rg, accountName, dbName)

		if !reflect.DeepEqual(err.Error(), "collection name is required") {
			t.Errorf("expected err %v", err)
		}
	})

	t.Run("error - resource group name is required", func(t *testing.T) {
		action := Create
		rg := ""

		_, err := ac.CreateCommand(action, rg, accountName, dbName)

		if !reflect.DeepEqual(err.Error(), "resource group name is required") {
			t.Errorf("expected err %v", err)
		}
	})

	t.Run("error - account name is required", func(t *testing.T) {
		action := Create
		accountName := ""

		_, err := ac.CreateCommand(action, rg, accountName, dbName)

		if !reflect.DeepEqual(err.Error(), "account name is required") {
			t.Errorf("expected err %v", err)
		}
	})

	t.Run("error - throughput and maxThroughput cannot be specified at the same time", func(t *testing.T) {
		action := Create

		_, err := ac.CreateCommand(action, rg, accountName, dbName)

		if !reflect.DeepEqual(err.Error(), "throughput and maxThroughput cannot be specified at the same time") {
			t.Errorf("expected err %v", err)
		}
	})

	t.Run("error - maxThroughput must be greater than or equal to 4000", func(t *testing.T) {
		action := Create
		ac.Throughput = new(int)
		*ac.Throughput = 100
		ac.AutoScale = new(bool)
		*ac.AutoScale = true
		ac.SharedRU = false

		_, err := ac.CreateCommand(action, rg, accountName, dbName)

		if !reflect.DeepEqual(err.Error(), "maxThroughput must be greater than or equal to 4000") {
			t.Errorf("expected err %v", err)
		}
	})

	t.Run("error - throughput must be greater than or equal to 400", func(t *testing.T) {
		action := Create
		ac.Throughput = new(int)
		*ac.Throughput = 100
		ac.AutoScale = new(bool)
		*ac.AutoScale = false
		ac.SharedRU = false

		_, err := ac.CreateCommand(action, rg, accountName, dbName)

		if !reflect.DeepEqual(err.Error(), "throughput must be greater than or equal to 400") {
			t.Errorf("expected err %v", err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		action := Delete
		expectedArgs := []string{"cosmosdb", "mongodb", "collection", "delete", "-g", rg, "-a", accountName, "-d", dbName, "-n", ac.Collection}

		args, _ := ac.CreateCommand(action, rg, accountName, dbName)

		if !reflect.DeepEqual(args, expectedArgs) {
			t.Errorf("expected args %v, got %v", expectedArgs, args)
		}
	})

	t.Run("List", func(t *testing.T) {
		action := List
		expectedArgs := []string{"cosmosdb", "mongodb", "collection", "list", "-g", rg, "-a", accountName, "-d", dbName}

		args, _ := ac.CreateCommand(action, rg, accountName, dbName)

		if !reflect.DeepEqual(args, expectedArgs) {
			t.Errorf("expected args %v, got %v", expectedArgs, args)
		}
	})

	t.Run("Show", func(t *testing.T) {
		action := Show
		expectedArgs := []string{"cosmosdb", "mongodb", "collection", "show", "-g", rg, "-a", accountName, "-d", dbName, "-n", ac.Collection}

		args, _ := ac.CreateCommand(action, rg, accountName, dbName)

		if !reflect.DeepEqual(args, expectedArgs) {
			t.Errorf("expected args %v, got %v", expectedArgs, args)
		}
	})

	t.Run("Exists", func(t *testing.T) {
		action := Exists
		expectedArgs := []string{"cosmosdb", "mongodb", "collection", "exists", "-g", rg, "-a", accountName, "-d", dbName, "-n", ac.Collection}

		args, _ := ac.CreateCommand(action, rg, accountName, dbName)

		if !reflect.DeepEqual(args, expectedArgs) {
			t.Errorf("expected args %v, got %v", expectedArgs, args)
		}
	})

	t.Run("error - invalid action", func(t *testing.T) {
		action := Action("invalid")

		_, err := ac.CreateCommand(action, rg, accountName, dbName)

		if !reflect.DeepEqual(err.Error(), "invalid action") {
			t.Errorf("expected err %v", err)
		}
	})
}

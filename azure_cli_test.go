package main

import (
	"reflect"
	"testing"
)

func TestCollection(t *testing.T) {
	ac := AzureCommand{
		Collection: "MyCollection",
		ShardKey:   nil,
		SharedRU:   false,
		AutoScale:  nil,
		Throughput: nil,
	}

	rg := "MyResourceGroup"
	accountName := "MyAccount"
	dbName := "MyDatabase"

	t.Run("Create", func(t *testing.T) {
		action := Create
		expectedArgs := []string{"cosmosdb", "mongodb", "collection", "create", "-g", rg, "-a", accountName, "-d", dbName, "-n", ac.Collection}

		args, _ := ac.CreateCommand(action, rg, accountName, dbName)

		if !reflect.DeepEqual(args, expectedArgs) {
			t.Errorf("expected args %v, got %v", expectedArgs, args)
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

	t.Run("InvalidAction", func(t *testing.T) {
		action := Action("invalid")
		expectedArgs := []string{}

		args, _ := ac.CreateCommand(action, rg, accountName, dbName)

		if !reflect.DeepEqual(args, expectedArgs) {
			t.Errorf("expected args %v, got %v", expectedArgs, args)
		}
	})
}

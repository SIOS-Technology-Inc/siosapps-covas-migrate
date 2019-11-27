package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var db *mongo.Database

// Handler returns database handler.
func handler() *mongo.Database {
	if db != nil {
		return db
	}

	u, err := ParseURI(os.Getenv("URI"))

	if err != nil {
		panic(fmt.Sprintf("incorrect URI given, %s", err))
	}

	client, err := mongo.Connect(ctx(), options.Client().ApplyURI(os.Getenv("URI")))

	if err != nil {
		panic(err)
	}

	if err := client.Ping(ctx(), readpref.Primary()); err != nil {
		panic(err)
	}

	db = client.Database(u.Database)

	return db
}

func ctx() context.Context {
	c, _ := context.WithTimeout(context.Background(), 30*time.Second)

	return c
}

const migrationCollection = "migrations"
const migrationKey = "latest"
const migrationInitValue = "0"

// Setup collection to store migration history.
func Setup() error {
	func() {
		var result bson.M

		q := bson.D{{migrationKey, bson.D{{"$exists", true}}}}
		opts := options.FindOne()

		err := handler().Collection(migrationCollection).FindOne(ctx(), q, opts).Decode(&result)

		if err != nil && err != mongo.ErrNoDocuments {
			panic(fmt.Sprintf("failed to lookup migration key, error %s", err))
		}

		if result != nil {
			panic(fmt.Sprintf("record already exists for migration key %s", migrationKey))
		}
	}()

	rec := bson.M{migrationKey: migrationInitValue}

	if _, err := handler().Collection(migrationCollection).InsertOne(ctx(), rec); err != nil {
		return fmt.Errorf("failed to setup initial collection, %s", err)
	}

	return nil
}

func filename(given string) string {
	chunks := strings.Split(given, "/")

	return chunks[len(chunks)-1]
}

// Current retrieves migration history.
func Current() (string, error) {
	q := bson.D{{migrationKey, bson.D{{"$exists", true}}}}
	opts := options.FindOne()

	var result bson.M

	if err := handler().Collection(migrationCollection).FindOne(ctx(), q, opts).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to retrieve migration history")
	}

	return result["latest"].(string), nil
}

// Next returns a migration target within the given directory.
func Next(dir, current string) (*Command, error) {

	paths, err := filepath.Glob(fmt.Sprintf("%s/*.json", dir))

	if err != nil {
		return nil, fmt.Errorf("failed to glob")
	}

	if len(paths) == 0 {
		return nil, fmt.Errorf("directory does not contain any schema files in JSON")
	}

	// Return the first match when called after init.
	if current == migrationInitValue {
		t := paths[0]

		cmd, err := parseCommand(t, handler().Name())

		if err != nil || cmd == nil {
			return nil, fmt.Errorf("failed to parse JSON, schema is possibly broken, %s", err)
		}

		return cmd, nil
	}

	for idx, p := range paths {
		// Break when loop reaches to the last element.
		if len(paths) == idx+1 {
			break
		}

		f := filename(p)

		if current != f {
			continue
		}

		// Matched to current item, attempt to get next one.
		cmd, err := parseCommand(paths[idx+1], handler().Name())

		if err != nil || cmd == nil {
			return nil, fmt.Errorf("failed to parse JSON, schema is possibly broken, %s", err)
		}

		return cmd, nil
	}

	return nil, nil
}

// Apply changes to target database.
func Apply(in *Command) error {
	if in == nil {
		return fmt.Errorf("invalid command given")
	}

	// Run admin command (optional)
	if in.Admin != "" {
		if err := func() error {
			var cmd bson.D

			if err := bson.UnmarshalExtJSON([]byte(in.Admin), true, &cmd); err != nil {
				return err
			}

			opts := options.RunCmd().SetReadPreference(readpref.Primary())

			var out bson.M

			// if err := handler().Client().Database("admin").RunCommand(ctx(), debug, opts).Decode(&out); err != nil {
			if err := handler().Client().Database("admin").RunCommand(ctx(), cmd, opts).Decode(&out); err != nil {
				return err
			}

			return nil
		}(); err != nil {
			return err
		}
	}

	// Run user command (optional)
	if in.General != "" {
		if err := func() error {
			var cmd bson.D

			if err := bson.UnmarshalExtJSON([]byte(in.General), true, &cmd); err != nil {
				return err
			}

			opts := options.RunCmd().SetReadPreference(readpref.Primary())

			var out bson.M

			if err := handler().RunCommand(ctx(), cmd, opts).Decode(&out); err != nil {
				return err
			}
			return nil
		}(); err != nil {
			return err
		}
	}

	// After everything is done, update state to be the latest.
	if err := func() error {
		opts := options.FindOneAndUpdate().SetUpsert(true)
		q := bson.D{{migrationKey, bson.D{{"$exists", true}}}}
		update := bson.D{{"$set", bson.D{{migrationKey, in.Version}}}}

		var updated bson.M

		return handler().Collection(migrationCollection).FindOneAndUpdate(ctx(), q, update, opts).Decode(&updated)
	}(); err != nil {
		return err
	}

	return nil
}

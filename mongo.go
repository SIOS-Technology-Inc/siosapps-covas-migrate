package main

import (
	"context"
	"encoding/json"
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

/*
Setup collection to store migration history.

マイグレーション履歴を保存するためのコレクションをセットアップします。
*/
func Setup() error {
	func() {
		var result bson.M

		q := bson.D{{Key: migrationKey, Value: bson.D{{Key: "$exists", Value: true}}}}
		opts := options.FindOne()

		err := handler().Collection(migrationCollection).FindOne(ctx(), q, opts).Decode(&result)

		if err != nil && err != mongo.ErrNoDocuments {
			panic(fmt.Sprintf("failed to lookup migration key, error %s", err))
		}

		if result != nil {
			panic(fmt.Sprintf("record already exists for migration key %s", migrationKey))
		}
	}()

	rec := bson.D{{Key: migrationKey, Value: migrationInitValue}}

	if _, err := handler().Collection(migrationCollection).InsertOne(ctx(), rec); err != nil {
		return fmt.Errorf("failed to setup initial collection, %s", err)
	}

	return nil
}

func filename(given string) string {
	chunks := strings.Split(given, "/")

	return chunks[len(chunks)-1]
}

/*
Current retrieves migration history.

マイグレーション履歴を取得します。
*/
func Current() (string, error) {
	q := bson.D{{Key: migrationKey, Value: bson.D{{Key: "$exists", Value: true}}}}
	opts := options.FindOne()

	var result bson.M

	if err := handler().Collection(migrationCollection).FindOne(ctx(), q, opts).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to retrieve migration history")
	}

	return result["latest"].(string), nil
}

/*
Next returns a migration target within the given directory.

指定されたディレクトリ内のマイグレーション対象を返します。
*/
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

/*
Apply changes to target database.

データベースに変更を適用します。

	@param in *Command
	@param u URI
*/
func Apply(in *Command, u *URI, rg string) error {
	if in == nil {
		return fmt.Errorf("invalid command given")
	}

	// Run admin command (optional)
	// TODO: Azure CLIに変更
	if in.Admin != "" {
		// ローカル環境用の処理
		if strings.Contains(u.Host, "localhost") {
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
		} else {
			// Azure環境用の処理
			if err := func() error {
				var cmd AzureCommand

				if err := json.Unmarshal([]byte(in.Admin), &cmd); err != nil {
					return err
				}
				fmt.Println(cmd.Description)
				opts, err := cmd.CreateCommand(Create, rg, u.Username, u.Database)
				if err != nil {
					return err
				}
				fmt.Println(opts)
				if _, err := AzExcute(opts); err != nil {
					return err
				}

				return nil
			}(); err != nil {
				return err
			}
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
		q := bson.D{{Key: migrationKey, Value: bson.D{{Key: "$exists", Value: true}}}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: migrationKey, Value: in.Version}}}}

		var updated bson.M

		return handler().Collection(migrationCollection).FindOneAndUpdate(ctx(), q, update, opts).Decode(&updated)
	}(); err != nil {
		return err
	}

	return nil
}

func Update(dirName, adminFlag string, u *URI, rg string) error {
	// Matched to current item, attempt to get next one.
	in, err := parseCommand(dirName, handler().Name())
	if err != nil || in == nil {
		return fmt.Errorf("failed to parse JSON, schema is possibly broken, %s", err)
	}

	// Run admin command (optional)
	if in.Admin != "" {
		if adminFlag != "true" {
			// ローカル環境用の処理
			if strings.Contains(u.Host, "localhost") {
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
			} else {
				// Azure環境用の処理
				if err := func() error {
					var cmd AzureCommand

					if err := json.Unmarshal([]byte(in.Admin), &cmd); err != nil {
						return err
					}
					fmt.Println(cmd.Description)
					opts, err := cmd.CreateCommand(Create, rg, u.Username, u.Database)
					if err != nil {
						return err
					}
					if _, err := AzExcute(opts); err != nil {
						return err
					}

					return nil
				}(); err != nil {
					return err
				}
			}
		}
	}

	// Run user command (optional)
	if in.General != "" {
		if err := func() error {
			var cmd bson.D

			if err := bson.UnmarshalExtJSON([]byte(in.General), true, &cmd); err != nil {
				return err
			}
			fmt.Println(cmd)

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

	return nil
}

func Revert(fileName string) error {
	if err := func() error {
		opts := options.FindOneAndUpdate().SetUpsert(true)
		q := bson.D{{Key: migrationKey, Value: bson.D{{Key: "$exists", Value: true}}}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: migrationKey, Value: fileName}}}}

		var updated bson.M

		return handler().Collection(migrationCollection).FindOneAndUpdate(ctx(), q, update, opts).Decode(&updated)
	}(); err != nil {
		return err
	}

	return nil
}

func FindIndex(collectionName string) error {
	indexView := handler().Collection(collectionName).Indexes()
	opts := options.ListIndexes().SetMaxTime(2 * time.Second)
	cursor, err := indexView.List(context.TODO(), opts)
	if err != nil {
		return err
	}
	var result []bson.M
	if err = cursor.All(context.TODO(), &result); err != nil {
		return err
	}
	for _, v := range result {
		for k1, v1 := range v {
			fmt.Printf("%v: %v\n", k1, v1)
		}
		fmt.Println()
	}
	return nil
}

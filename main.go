package main

import (
	"fmt"
	"os"
	"time"

	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "migrate",
		Usage:   "MongoDB migration tool with minimal api",
		Version: "0.7.0",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "Setup",
				Action: func(c *cli.Context) error {
					out := Setup()
					fmt.Println(out)
					return out
				},
			},
			{
				Name:  "up",
				Usage: "Run migration",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "dir",
						Aliases: []string{"d"},
						Value:   "migrations",
						Usage:   "Directory of your JSON files",
					},
					&cli.StringFlag{
						Name:    "rg",
						Aliases: []string{"r"},
						Value:   "MyResourceGroup",
						Usage:   "Resource Group Name",
					},
				},
				Action: func(c *cli.Context) error {
					for {
						fmt.Printf("checking current state.. ")
						cur, err := Current()

						if err != nil {
							fmt.Printf("failed, %s \n", err)
							return err
						}

						fmt.Printf("ok \n")

						fmt.Printf("retrieving the changes.. ")

						next, err := Next(c.String("dir"), cur)

						if err != nil {
							fmt.Printf("failed, %s \n", err)
							return err
						}

						if next == nil {
							fmt.Printf("no more migrations. \n")
							break
						}

						fmt.Printf("%s \n", next.Version)
						fmt.Printf("applying changes.. ")

						resourceGroup := c.String("rg")
						u, _ := ParseURI(os.Getenv("URI"))
						if err := Apply(next, u, resourceGroup); err != nil {
							fmt.Printf("failed, %s \n", err)
							return err
						}

						fmt.Printf("completed migration. \n")

						// Set interval to reduce database load.
						time.Sleep(2 * time.Second)
					}

					fmt.Println("done!")
					return nil
				},
			},
			{
				Name:  "fix",
				Usage: "Run migration",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Value:   "migrations",
						Usage:   "Input your JSON file",
					},
					&cli.StringFlag{
						Name:    "admin",
						Aliases: []string{"a"},
						Value:   "false",
						Usage:   "AdminCommand is run, if flag is true.",
					},
					&cli.StringFlag{
						Name:    "rg",
						Aliases: []string{"r"},
						Value:   "MyResourceGroup",
						Usage:   "Resource Group Name",
					},
				},
				Action: func(c *cli.Context) error {
					file := c.String("file")
					adminFlag := c.String("admin")
					resourceGroup := c.String("rg")
					u, _ := ParseURI(os.Getenv("URI"))
					if err := Update(file, adminFlag, u, resourceGroup); err != nil {
						fmt.Printf("failed, %s \n", err)
						return err
					}

					fmt.Println(file)
					fmt.Println("done!")
					return nil
				},
			},
			{
				Name:  "revert",
				Usage: "Revert migration pointer",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "pointer file name",
						Usage:   "Target file you want to revert",
					},
				},
				Action: func(c *cli.Context) error {
					fileName := c.String("name")
					if err := Revert(fileName); err != nil {
						fmt.Printf("failed, %s \n", err)
						return err
					}

					fmt.Println(fileName)
					fmt.Println("done!")
					return nil
				},
			},
			{
				Name:  "index",
				Usage: "find index by collection name",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "",
						Usage:   "Target collection you want to index",
					},
				},
				Action: func(c *cli.Context) error {
					collectionName := c.String("name")

					if collectionName == "" {
						fmt.Println("required name option")
						return nil
					}

					if err := FindIndex(collectionName); err != nil {
						fmt.Printf("failed, %s \n", err)
						return err
					}

					fmt.Println(collectionName)
					fmt.Println("done!")
					return nil
				},
			},
			{
				Name:  "delete",
				Usage: "delete index by collection name and index name",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "collection",
						Aliases: []string{"c"},
						Value:   "",
						Usage:   "Target collection you want to delete index",
					},
					&cli.StringFlag{
						Name:    "index",
						Aliases: []string{"i"},
						Value:   "",
						Usage:   "Target index you want to delete",
					},
				},
				Action: func(c *cli.Context) error {
					collectionName := c.String("collection")
					indexName := c.String("index")

					if collectionName == "" || indexName == "" {
						fmt.Println("required collection and index option")
						return nil
					}

					if err := DeleteIndex(collectionName, indexName); err != nil {
						fmt.Printf("failed, %s \n", err)
						return err
					}

					fmt.Printf("collection: %s, index: %s \n", collectionName, indexName)
					fmt.Println("done!")
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		// エラーがある場合は異常終了にしたい
		os.Exit(1)
	}
}

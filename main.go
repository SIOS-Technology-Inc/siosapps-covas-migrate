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
		Version: "0.5.0",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "Setup",
				Action: func(c *cli.Context) error {
					return Setup()
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

						if err := Apply(next); err != nil {
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
				},
				Action: func(c *cli.Context) error {
					file := c.String("file")
					adminFlag := c.String("admin")
					if err := Update(file, adminFlag); err != nil {
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
		},
	}

	app.Run(os.Args)
}

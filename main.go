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
		Version: "0.1.0",
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

						fmt.Printf("%s \n", next.Version)

						if next == nil {
							fmt.Printf("no more migrations\n")
							break
						}

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
		},
	}

	app.Run(os.Args)
}

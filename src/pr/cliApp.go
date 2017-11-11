package pr

import (
	"context"
	"os"
	"poprep/src/prconf"
	"poprep/src/prgithub"

	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

var client = &github.Client{}
var ctx context.Context

func initCliApp() *cli.App {
	ctx, client = prgithub.Get()
	app := cli.NewApp()
	app.Name = prconf.AppName
	app.Usage = prconf.AppUsage
	app.Version = prconf.AppVersion
	app.Compiled = prconf.AppCompiled
	app.Authors = prconf.AppAuthors
	app.Copyright = prconf.AppCopyright

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "lists all available awesome-lists",
			Action: func(c *cli.Context) error {
				if !c.Args().Present() {
					table := tablewriter.NewWriter(os.Stdout)
					table.SetHeader([]string{"Command:", "Name:", "URL:", "Created by:"})
					for k, v := range config.AwesomeLists {
						table.Append([]string{"list " + k, v.Name, v.URL, v.Author})
					}
					table.Render()
				}
				return nil
			},
			Subcommands: getListSubcommands(),
		},
		{
			Name:    "categories",
			Aliases: []string{"cat"},
			Usage:   "lists all available categories of specified language",
			Action: func(c *cli.Context) error {
				if !c.Args().Present() {
					cli.ShowAppHelp(c)
				}
				return nil
			},
			Subcommands: getCatSubcommands(),
		},
		{
			Name:    "cache",
			Aliases: []string{"cc"},
			Usage:   "creates a local cache so the next commands run faster",
			Action: func(c *cli.Context) error {
				if !c.Args().Present() {
					cli.ShowAppHelp(c)
				}
				return nil
			},
			Subcommands: getCacheSubcommands(),
		},
	}
	return app
}

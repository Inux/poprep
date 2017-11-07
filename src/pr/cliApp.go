package pr

import (
	"context"
	"fmt"
	"os"
	"poprep/src/prconf"
	"poprep/src/prcrawler"
	"poprep/src/prgithub"
	"strconv"

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
	}
	return app
}

//getListSubcommands returns the poprep list cli commands
func getListSubcommands() []cli.Command {
	subcommands := []cli.Command{}

	//Add subcommands to display all available repositories
	for k, v := range config.AwesomeLists {
		cmd := new(cli.Command)
		cmd.Name = k
		cmd.Usage = "list " + k
		cmd.UsageText = "list " + k + "(for " + v.Name + ")"
		cmd.Action = actionPrintAwesomeList
		subcommands = append(subcommands, *cmd)
	}
	return subcommands
}

func actionPrintAwesomeList(c *cli.Context) error {
	if !c.Args().Present() {
		ri := prcrawler.GetRepoInfos(ctx, client,
			config.AwesomeLists[c.Command.Name], "")

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name:", "Category:", "Stars:", "URL:", "Created by:"})
		for _, v := range ri {
			table.Append([]string{v.Project, v.Category, strconv.Itoa(v.Stars), v.URL, v.User})
		}
		table.Render()
	} else {
		cli.ShowAppHelp(c)
	}
	return nil
}

//Subcommands for the "poprep" cat command
func getCatSubcommands() []cli.Command {
	subcommands := []cli.Command{}

	//Add subcommands to display all available categories
	for k, v := range config.AwesomeLists {
		cmd := new(cli.Command)
		cmd.Name = k
		cmd.Usage = "categories " + k
		cmd.UsageText = "categories " + k + "(for " + v.Name + ")"
		cmd.Action = actionPrintCategories
		subcommands = append(subcommands, *cmd)
	}
	return subcommands
}

func actionPrintCategories(c *cli.Context) error {
	if !c.Args().Present() {
		categories := prcrawler.GetCategories(ctx, config.AwesomeLists[c.Command.Name])

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Categories:"})
		for _, v := range categories {
			table.Append([]string{v})
		}
		table.Render()
	} else {
		if c.NArg() == 1 {
			language := c.Command.Name
			category := c.Args().Get(0)
			awesomeList := config.AwesomeLists[language]

			if prcrawler.ContainsCategory(awesomeList, category) {
				ri := prcrawler.GetRepoInfos(ctx, client,
					awesomeList, category)

				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Name:", "Category:", "Stars:", "URL:", "Created by:"})
				for _, v := range ri {
					table.Append([]string{v.Project, v.Category, strconv.Itoa(v.Stars), v.URL, v.User})
				}
				table.Render()
			} else {
				fmt.Println("Error:", awesomeList.Name, "doesn't contain category:", category)
			}
		}
	}
	return nil
}

//Package pr - Contains the category (cat) commands and actions
package pr

import (
	"fmt"
	"os"
	"poprep/src/prcrawler"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

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

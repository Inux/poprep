//Package pr - Contains the list (ls) commands and actions
package pr

import (
	"os"
	"poprep/src/prcrawler"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

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

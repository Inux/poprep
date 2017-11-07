package main

import (
	"poprep/src/pr"
)

const (
	awesomeListConfig = "poprep_config.toml"
	githubConfig      = "github_config.toml"
)

func main() {
	//init and run poprep
	pr.Init(awesomeListConfig, githubConfig)
	pr.RunCli()
}
